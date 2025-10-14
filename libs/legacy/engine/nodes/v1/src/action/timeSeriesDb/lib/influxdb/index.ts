import { Thing } from "@prisma/client"

export type InfluxDBNodeConfig = {
  "influxdbDatabase::thing"?: string | Thing
  bucket?: string
}

