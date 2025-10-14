package com.iotea.nodes.action.log

import com.iotea.nodes.api.Node
import com.iotea.nodes.api.NodeInfo
import com.iotea.nodes.api.HostContext
import com.google.auto.service.AutoService
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.launch
import javax.lang.model.type.NullType

@AutoService(Node::class)
public class LogActionNode : Node<ActionLogV01_0, ByteArray, NullType> {
  override val info = NodeInfo(
    abiVersion = 1, 
    nodeVersion = "0.1.0", 
    name = "action-log", 
    label = "Log", 
    description = "Logs a message to the console", 
    category = "action", 
    inputs = emptyList(), 
    outputs = emptyList()
  )
  
  override lateinit var hostCtx: HostContext
  override lateinit var scope: CoroutineScope
  override lateinit var inputChannel: Channel<ByteArray>
  override lateinit var outputChannels: List<Channel<NullType>>
  override lateinit var config: ActionLogV01_0

  override suspend fun start(
    hostCtx: HostContext,
    scope: CoroutineScope,
    config: ActionLogV01_0,
    inputChannel: Channel<ByteArray>,
    outputChannels: List<Channel<NullType>>
  ) {
    super.start(hostCtx, scope, config, inputChannel, emptyList())

    scope.launch {
      for (input in inputChannel) {
        handle(input)
      }
    }
  }

  override suspend fun handle(input: ByteArray) {
    println(String(input, Charsets.UTF_8))
  }

  override suspend fun stop() {
    super.stop()
  }
}