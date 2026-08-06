package com.gruent.nodes.api

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.cancel
import kotlinx.coroutines.channels.Channel
import java.nio.ByteBuffer

data class NodeInputOutput(val id: String, val label: String)

data class NodeInfo(
  val abiVersion: Int,
  val nodeVersion: String,
  val name: String,
  val label: String,
  val description: String,
  val category: String,
  val inputs: List<NodeInputOutput>,
  val outputs: List<NodeInputOutput>
)

interface HostContext {
  val nowNanos: Long
  // add logging/metrics helpers later
}

// Generic Node interface with type parameter for configuration
interface Node<Config> {
  val info: NodeInfo
  var hostCtx: HostContext
  var scope: CoroutineScope
  var config: Config
  var inputChannel: Channel<ByteBuffer>
  var outputChannels: List<Channel<ByteBuffer>>

  suspend fun start(
    hostCtx: HostContext,
    scope: CoroutineScope,
    config: Config,
    inputChannel: Channel<ByteBuffer>,
    outputChannels: List<Channel<ByteBuffer>>
  ) {
    this.hostCtx = hostCtx
    this.scope = scope
    this.config = config
    this.inputChannel = inputChannel
    this.outputChannels = outputChannels
  }

  suspend fun stop() {
    // Close input and output channels
    inputChannel.close()
    for (channel in outputChannels) {
      channel.close()
    }

    // Shut down the coroutine scope
    scope.cancel()
  }
}
