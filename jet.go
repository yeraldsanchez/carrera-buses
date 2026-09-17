package main

var jetGrid = []string{
	"..........TT...............",
	"...........TT..............",
	"............TT.............",
	"............TTT............",
	"...T........TTTT...........",
	"....T......TTTTTT..........",
	"...TTT....TTTTTTTTT........",
	"..TTTTT..TTTTWWTTTTT.......",
	".TTTTTTTTTTTWWWTTTTTTTT....",
	"DDFFFFFFFFFFFFFFFFFFFFNYNN.",
	"..DDFFFFFFFFFFFFFFFFFFNYYNN",
	"DDFFFFFFFFFFFFFFFFFFFFNYNN.",
	".TTTTTTTTTTTWWWTTTTTTTT....",
	"..TTTTT..TTTTWWTTTTT.......",
	"...TTT....TTTTTTTTT........",
	"....T......TTTTTT..........",
	"...T........TTTT...........",
	"............TTT............",
	"............TT.............",
	"...........TT..............",
	"..........TT...............",
}

var jetPalette = map[byte]string{
	'D': "#1A237E", // cola / timón
	'F': "#B0BEC5", // fuselaje
	'Y': "#4FC3F7", // cabina (cristal)
	'N': "#263238", // punta de la nariz
	'W': "#78909C", // ala
	'T': "#FF3B30", // borde/punta de ala (acento)
}
