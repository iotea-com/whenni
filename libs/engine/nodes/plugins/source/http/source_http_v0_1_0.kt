// To parse the JSON, install Klaxon and do:
//
//   val sourceHTTPV010 = SourceHTTPV01_0.fromJson(jsonString)

package quicktype

import com.beust.klaxon.*

private fun <T> Klaxon.convert(k: kotlin.reflect.KClass<*>, fromJson: (JsonValue) -> T, toJson: (T) -> String, isUnion: Boolean = false) =
    this.converter(object: Converter {
        @Suppress("UNCHECKED_CAST")
        override fun toJson(value: Any)        = toJson(value as T)
        override fun fromJson(jv: JsonValue)   = fromJson(jv) as Any
        override fun canConvert(cls: Class<*>) = cls == k.java || (isUnion && cls.superclass == k.java)
    })

private val klaxon = Klaxon()
    .convert(Method::class, { Method.fromValue(it.string!!) }, { "\"${it.value}\"" })

data class SourceHTTPV01_0 (
    val method: Method? = null,

    /**
     * The timeout for the response in milliseconds
     */
    val responseTimeout: Double? = null
) {
    public fun toJson() = klaxon.toJsonString(this)

    companion object {
        public fun fromJson(json: String) = klaxon.parse<SourceHTTPV01_0>(json)
    }
}

enum class Method(val value: String) {
    Delete("DELETE"),
    Get("GET"),
    Patch("PATCH"),
    Post("POST"),
    Put("PUT");

    companion object {
        public fun fromValue(value: String): Method = when (value) {
            "DELETE" -> Delete
            "GET"    -> Get
            "PATCH"  -> Patch
            "POST"   -> Post
            "PUT"    -> Put
            else     -> throw IllegalArgumentException()
        }
    }
}
