type Point = {
  x: number
  y: number
}

export const isPointInPolygon = (point: Point, polygon: Point[], canvasScale: number): boolean => {
  let x = point.x * canvasScale,
    y = point.y * canvasScale

  let inside = false
  for (let i = 0, j = polygon.length - 1; i < polygon.length; j = i++) {
    let xi = polygon[i].x,
      yi = polygon[i].y
    let xj = polygon[j].x,
      yj = polygon[j].y

    let intersect = yi > y != yj > y && x < ((xj - xi) * (y - yi)) / (yj - yi) + xi
    if (intersect) inside = !inside
  }

  return inside
}
