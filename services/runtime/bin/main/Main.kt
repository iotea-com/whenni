package com.gruent.runtime

import com.gruent.nodes.api.Node
import com.gruent.nodes.api.HostContext
import java.io.FileNotFoundException
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.Channel
import java.net.URLClassLoader
import java.nio.file.Path
import java.nio.file.Paths
import java.util.*

object HostCtx : HostContext {
  override val nowNanos: Long get() = System.nanoTime()
}

class PluginLoader {
  fun loadPlugin(jarPath: Path): List<Node<*, *, *>> {
    println("Loading node plugin from $jarPath...")
    val nodeFile = jarPath.toFile()

    // Check if the jar exists and throw an error if it doesn't
    if (!nodeFile.exists()) {
      throw FileNotFoundException("JAR file not found: $jarPath")
    }

    val loader = URLClassLoader(
      arrayOf(jarPath.toUri().toURL()),
      Node::class.java.classLoader
    )

    val serviceConfigUrl = loader.getResource("META-INF/services/com.gruent.nodes.api.Node")
    if (serviceConfigUrl != null) {
      val services = serviceConfigUrl.readText()
      println("Services registered from node plugin: $services")
    }

    // Use ServiceLoader to discover all Node implementations in the JAR
    val serviceLoader = ServiceLoader.load(Node::class.java, loader)
    val nodes = serviceLoader.toList()
    return nodes
  }

  fun loadPluginsFromDirectory(pluginDir: Path): List<Node<*, *, *>> {
    val plugins = mutableListOf<Node<*, *, *>>()

    pluginDir.toFile().listFiles { file ->
      file.isFile && file.name.endsWith(".jar")
    }?.forEach { jarFile ->
      try {
        val loadedPlugins = loadPlugin(jarFile.toPath())
        plugins.addAll(loadedPlugins)
        println("Loaded ${loadedPlugins.size} nodes from ${jarFile.name}")
      } catch (e: Exception) {
        println("Failed to load plugin from ${jarFile.name}: ${e.message}")
      }
    }

    return plugins
  }
}

@OptIn(DelicateCoroutinesApi::class)
fun main() {
  println("Runtime starting…")

  val pluginLoader = PluginLoader()

  val pluginDirPath = System.getenv("PLUGIN_DIR")?.let { Paths.get(it) }
    ?: Paths.get("build/dist/nodes")

  if (pluginDirPath.toFile().exists()) {
    val plugins = pluginLoader.loadPluginsFromDirectory(pluginDirPath)

    println("Plugins loaded: ${plugins.size} nodes")

    plugins.firstOrNull()?.let { node ->
      runBlocking {
        val scope = CoroutineScope(Dispatchers.Default)
        val inputChannel = Channel<ByteArray>()
        val outputChannels = emptyList<Channel<Any>>()

        try {
          // Create appropriate config based on node type
          when (node.info.name) {
            "action-log" -> {
              // For demonstration purposes, create a simple JSON config.
              // In a real channel, this config would be sourced from the database and compiled into a plan by the
              // orchestrator.
              val configJson = """{"message": "Hello from runtime!", "templateModel": null}"""

              // For the log action node, create ActionLogV01_0 config directly
              val configClass = node.javaClass.classLoader
                .loadClass("com.gruent.nodes.action.log.ActionLogV01_0")
              val constructor = configClass.getConstructor(String::class.java, String::class.java)
              val config = constructor.newInstance("Hello from runtime!", null)

              // Cast to the generic Node interface and call start directly
              @Suppress("UNCHECKED_CAST")
              val typedNode = node as com.gruent.nodes.api.Node<Any, ByteArray, Any>
              typedNode.start(HostCtx, scope, config, inputChannel, outputChannels)

              // Start a co-routine that publishes messages
              val doneChannel = Channel<Boolean>()

              val sendInputJob = scope.launch {
                repeat(10) {
                  val now = System.currentTimeMillis().toString()
                  inputChannel.send("current time: $now".toByteArray())
                  delay(1000)
                }
                doneChannel.send(true)
              }

              println("awaiting done message...")
              doneChannel.receive()
              println("done!")
            }

            else -> {
              // For other nodes, use a generic approach
              @Suppress("UNCHECKED_CAST")
              val genericNode = node as com.gruent.nodes.api.Node<Any, ByteArray, Any>
              genericNode.start(HostCtx, scope, Any(), inputChannel, outputChannels)
            }
          }

          // Clean shutdown
          delay(1000)
          node.stop()
        } catch (e: Exception) {
          println("Error running node: ${e.message}")
          e.printStackTrace()
        }
      }
    }
  } else {
    println("Plugin directory not found: $pluginDirPath")
    println("Build plugins first with: ./gradlew buildNodes")
  }
}