plugins {
    kotlin("jvm") version "2.0.0"
    application
}

java { toolchain { languageVersion.set(JavaLanguageVersion.of(21)) } }

dependencies {
    implementation(project(":libs:engine:nodes:api"))

    // Logging / metrics
    implementation("ch.qos.logback:logback-classic:1.5.6")
    implementation("io.micrometer:micrometer-registry-prometheus:1.13.1")

    // Optional HTTP, Kafka, etc., later
    // implementation("io.netty:netty-all:4.1.111.Final")
}

application {
    mainClass.set("com.iotea.runtime.MainKt")
}

tasks.withType<Jar> {
    manifest {
        attributes["Main-Class"] = "com.iotea.runtime.MainKt"
    }
}

tasks.test { useJUnitPlatform() }
