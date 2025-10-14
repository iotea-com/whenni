plugins {
    kotlin("jvm") version "2.0.0"
    id("com.google.devtools.ksp") version "2.0.0-1.0.24" // match your Kotlin
    `java-library`
}

java { toolchain { languageVersion.set(JavaLanguageVersion.of(21)) } }

dependencies {
    implementation(project(":libs:engine:nodes:api"))
    ksp("dev.zacsweers.autoservice:auto-service-ksp:1.2.0")
    compileOnly("com.google.auto.service:auto-service-annotations:1.1.1")
    implementation(libs.klaxon)
    // node-specific libs
}

tasks.test { useJUnitPlatform() }

// Thin "bundle" (no shading): copy jar + runtime deps into /build/node-dist/lib
tasks.register<Copy>("bundleNode") {
    dependsOn(configurations.runtimeClasspath, tasks.named("jar"))
    val outDir = layout.buildDirectory.dir("node-dist").get().asFile
    
    // Set the main destination directory
    into(File(outDir, "lib"))
    
    // Configure what to copy
    from(configurations.runtimeClasspath.get().filter { it.name.endsWith(".jar") })
    from(tasks.named("jar"))
    
    // IMPORTANT: do NOT ship Kotlin stdlib or coroutines here — host provides them
    exclude("kotlin-stdlib*.jar", "kotlin-stdlib-jdk*.jar", "kotlin-reflect*.jar")
    exclude("kotlinx-coroutines-*.jar")
    exclude("libs-engine-nodes-api-*.jar") // don't duplicate the SPI
    
    // Copy metadata to the parent directory
    from("node.json") { 
        into(outDir)
    }
}

// Read node metadata from node.json
val nodeJsonFile = file("node.json")
val nodeJson = groovy.json.JsonSlurper().parseText(nodeJsonFile.readText()) as Map<String, Any>
val nodeCategory = nodeJson["category"] as String
val nodeName = nodeJson["name"] as String  // Using "name" instead of "label" for filename
val nodeVersion = nodeJson["version"] as String  // Using "name" instead of "label" for filename

// Configure JAR naming based on node metadata
tasks.jar {
    val nodeReference = "$nodeCategory-$nodeName"
    val version = nodeVersion
    archiveBaseName.set("$nodeCategory-$nodeName")

    
    // Output jar size after the jar is built
    doLast {
        val jarFile = archiveFile.get().asFile
        if (jarFile.exists()) {
            val sizeInBytes = jarFile.length()
            val sizeInKB = sizeInBytes / 1024.0
            val sizeInMB = sizeInKB / 1024.0
            
            println("Completed build for node $nodeReference v$version")
            println("JAR size information:")
            println("  File: ${jarFile.name}")
            println("  Size: %.2f KB (%.2f MB)".format(sizeInKB, sizeInMB))
        } else {
            println("JAR file not found: ${jarFile.absolutePath}")
        }
    }
}

// Task to output jar size information
tasks.register("jarSize") {
    dependsOn(tasks.jar)
    doLast {
        val jarTask = tasks.jar.get()
        val jarFile = jarTask.archiveFile.get().asFile
        
        if (jarFile.exists()) {
            val sizeInBytes = jarFile.length()
            val sizeInKB = sizeInBytes / 1024.0
            val sizeInMB = sizeInKB / 1024.0
            
            println("═══════════════════════════════════")
            println("JAR SIZE REPORT")
            println("═══════════════════════════════════")
            println("File: ${jarFile.name}")
            println("Path: ${jarFile.absolutePath}")
            println("Size: $sizeInBytes bytes")
            println("Size: %.2f KB".format(sizeInKB))
            println("Size: %.2f MB".format(sizeInMB))
            println("═══════════════════════════════════")
        } else {
            println("❌ JAR file not found: ${jarFile.absolutePath}")
        }
    }
}

group = "com.iotea.nodes"
version = "0.1.0"