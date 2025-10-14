plugins {
  // make Kotlin available, but do NOT apply it here
  id("org.jetbrains.kotlin.jvm") version "2.0.0" apply false
}

allprojects {
  repositories { mavenCentral() }
}

subprojects {
  // Configure ONLY projects that actually apply the Kotlin plugin
  plugins.withId("org.jetbrains.kotlin.jvm") {
    // Toolchain
    extensions.configure(
      org.jetbrains.kotlin.gradle.dsl.KotlinJvmProjectExtension::class.java
    ) {
      jvmToolchain(21)
    }

    // Kotlin compile options (Kotlin 2.x API)
    tasks.withType(org.jetbrains.kotlin.gradle.tasks.KotlinCompile::class.java).configureEach {
      compilerOptions.jvmTarget.set(org.jetbrains.kotlin.gradle.dsl.JvmTarget.JVM_21)
      // compilerOptions.freeCompilerArgs.add("-Xjsr305=strict")
    }

    // Tests on JUnit Platform
    tasks.withType(Test::class.java).configureEach {
      useJUnitPlatform()
    }
  }

  // Java-only projects (like :libs:engine:nodes:api) still get JUnit Platform
  tasks.withType(Test::class.java).configureEach { useJUnitPlatform() }
}

// ---- Helper aggregate tasks
tasks.register("buildNodes") {
  // Extract just the task paths during configuration
  val bundleTaskPaths = subprojects
    .filter { p ->
      p.path.startsWith(":libs:engine:nodes:plugins:") &&
      p.tasks.findByName("bundleNode") != null
    }
    .map { "${it.path}:bundleNode" }
  
  dependsOn(bundleTaskPaths)
}

tasks.register("distNodes") {
  dependsOn("buildNodes")
  
  // Extract node metadata during configuration phase
  val nodeProjectsData = subprojects
    .filter { it.path.startsWith(":libs:engine:nodes:plugins:") }
    .mapNotNull { p ->
      val nodeJsonFile = p.projectDir.resolve("node.json")
      if (!nodeJsonFile.exists()) {
        null
      } else {
        try {
          // Read and parse JSON during configuration time
          val nodeJsonText = nodeJsonFile.readText()
          val nodeJson = groovy.json.JsonSlurper().parseText(nodeJsonText) as Map<*, *>
          mapOf(
            "path" to p.path,
            "projectDir" to p.projectDir.absolutePath,
            "category" to (nodeJson["category"] as String),
            "name" to (nodeJson["name"] as String),
            "version" to (nodeJson["version"] as String)
          )
        } catch (e: Exception) {
          println("Warning: Could not parse node.json for project ${p.path}: ${e.message}")
          null
        }
      }
    }
  
  // Capture the root layout during configuration
  val rootBuildDir = layout.buildDirectory
  
  doLast {
    val dest = rootBuildDir.dir("dist/nodes").get().asFile
    dest.mkdirs()

    nodeProjectsData.forEach { projectInfo ->
      val projectDir = projectInfo["projectDir"] as String
      val nodeCategory = projectInfo["category"] as String
      val nodeName = projectInfo["name"] as String
      val nodeVersion = projectInfo["version"] as String
      
      
      // Construct the jar path using standard File operations
      val jarFile = File("$projectDir/build/libs/$nodeCategory-$nodeName-$nodeVersion.jar")
      
      if (jarFile.exists()) {
        println("Copying node $nodeCategory-$nodeName-$nodeVersion to $dest/$nodeCategory-$nodeName-$nodeVersion.jar...")
        jarFile.copyTo(
          target = File(dest, "$nodeCategory-$nodeName-$nodeVersion.jar"),
          overwrite = true
        )
      } else {
        println("Warning: JAR file not found for $nodeCategory-$nodeName-$nodeVersion: $jarFile")
      }
    }
  }
}

tasks.register("buildRuntime") {
  // Extract serializable data during configuration
  val runtimeProjectsData = subprojects
    .filter { it.path.startsWith(":services:runtime") }
    .map { p ->
      mapOf(
        "name" to p.name,
        "path" to p.path
      )
    }
  
  // Get the root project layout during configuration
  val rootBuildDir = layout.buildDirectory
  
  doLast {
    val dest = rootBuildDir.dir("dist/runtime").get().asFile
    dest.mkdirs()
    
    runtimeProjectsData.forEach { projectInfo ->
      val projectName = projectInfo["name"] as String
      val projectPath = projectInfo["path"] as String
      
      // Construct the jar path using the captured root build dir
      val jarFile = rootBuildDir
        .dir("../..${projectPath.replace(':', '/')}/build/libs/${projectName}-all.jar")
        .get().asFile
      
      if (jarFile.exists()) {
        jarFile.copyTo(File(dest, jarFile.name), overwrite = true)
      }
    }
  }
}