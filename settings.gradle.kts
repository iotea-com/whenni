pluginManagement {
  repositories {
    gradlePluginPortal()
    mavenCentral()
  }
}

plugins {
  id("org.gradle.toolchains.foojay-resolver-convention") version "1.0.0"
}

rootProject.name = "iotea"

// ---- Explicit modules
include(":services:runtime")
include(":libs:engine:nodes:api")

// ---- Dynamic discovery of plugin modules
// Expect structure: libs/engine/nodes/plugins/<category>/<label>/<version>
val pluginsRoot = file("libs/engine/nodes/plugins")
if (pluginsRoot.exists()) {
  pluginsRoot.listFiles { f -> f.isDirectory }?.forEach { category ->
    category.listFiles { f -> f.isDirectory }?.forEach { label ->
      label.listFiles { f -> f.isDirectory }?.forEach { version ->
        val path = ":libs:engine:nodes:plugins:${category.name}:${label.name}:${version.name}"
        include(path)
        project(path).projectDir = version
      }
    }
  }
}
