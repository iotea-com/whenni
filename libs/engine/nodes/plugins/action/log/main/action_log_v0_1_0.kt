// To parse the JSON, install Klaxon and do:
//
//   val actionLogV010 = ActionLogV01_0.fromJson(jsonString)

package com.iotea.nodes.action.log

import com.beust.klaxon.*

private val klaxon = Klaxon()

data class ActionLogV01_0 (
    /**
     * The message to log
     */
    val message: String,

    /**
     * The ID of the model to use if the message includes dynamic content
     */
    val templateModel: String? = null
) {
    public fun toJson() = klaxon.toJsonString(this)

    companion object {
        public fun fromJson(json: String) = klaxon.parse<ActionLogV01_0>(json)
    }
}
