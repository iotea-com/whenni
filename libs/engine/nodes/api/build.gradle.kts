plugins {
    kotlin("jvm") version "2.0.0"
    `java-library`
}

repositories {
    // Use Maven Central for resolving dependencies.
    mavenCentral()
}

// Apply a specific Java toolchain to ease working on different environments.
java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}


dependencies {
    api("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.8.1")
}

tasks.test { useJUnitPlatform() }