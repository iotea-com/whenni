// To parse the JSON, install Klaxon and do:
//
//   val conditionalThresholdV010 = ConditionalThresholdV01_0.fromJson(jsonString)

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
    .convert(Operator::class,        { Operator.fromValue(it.string!!) },        { "\"${it.value}\"" })
    .convert(LogicalOperator::class, { LogicalOperator.fromValue(it.string!!) }, { "\"${it.value}\"" })

data class ConditionalThresholdV01_0 (
    /**
     * The threshold conditions to evaluate
     */
    val conditions: List<ConfigSchema>? = null,

    /**
     * The logical operator to use for evaluating multiple conditions
     */
    val logicalOperator: LogicalOperator
) {
    public fun toJson() = klaxon.toJsonString(this)

    companion object {
        public fun fromJson(json: String) = klaxon.parse<ConditionalThresholdV01_0>(json)
    }
}

data class ConfigSchema (
    /**
     * The ID of the attribute to evaluate
     */
    @Json(name = "attributeId")
    val attributeID: String? = null,

    /**
     * The operator to use for evaluating the condition
     */
    val operator: Operator? = null,

    /**
     * The value to compare against
     */
    val value: Double? = null
)

/**
 * The operator to use for evaluating the condition
 */
enum class Operator(val value: String) {
    EqualTo("equal-to"),
    GreaterThan("greater-than"),
    GreaterThanEqualTo("greater-than-equal-to"),
    LessThan("less-than"),
    LessThanEqualTo("less-than-equal-to"),
    NotEqualTo("not-equal-to");

    companion object {
        public fun fromValue(value: String): Operator = when (value) {
            "equal-to"              -> EqualTo
            "greater-than"          -> GreaterThan
            "greater-than-equal-to" -> GreaterThanEqualTo
            "less-than"             -> LessThan
            "less-than-equal-to"    -> LessThanEqualTo
            "not-equal-to"          -> NotEqualTo
            else                    -> throw IllegalArgumentException()
        }
    }
}

/**
 * The logical operator to use for evaluating multiple conditions
 */
enum class LogicalOperator(val value: String) {
    And("AND"),
    Or("OR");

    companion object {
        public fun fromValue(value: String): LogicalOperator = when (value) {
            "AND" -> And
            "OR"  -> Or
            else  -> throw IllegalArgumentException()
        }
    }
}
