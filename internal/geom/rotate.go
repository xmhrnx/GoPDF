package geom

// RotateCW returns p rotated a quarter turn clockwise on a page whose
// display height is h: the new page is h wide.
func RotateCW(p Point, h float64) Point { return Point{h - p.Y, p.X} }

// RotateCCW returns p rotated a quarter turn counter-clockwise on a page
// whose display width is w: the new page is w tall.
func RotateCCW(p Point, w float64) Point { return Point{p.Y, w - p.X} }
