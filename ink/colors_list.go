package ink

// Named color variables provide a convenient palette for use with [Style.WithForeground]
// and [Style.WithBackground]. All values are RGB true-color [Color] values.
//
// The palette is organised into thematic groups:
//   - Basic & bright ANSI-equivalent colors ([Black]–[Cyan], [BrightBlack]–[BrightCyan])
//   - Dark variants ([DarkRed], [DarkGreen], …)
//   - Grays & neutrals ([Gray], [Silver], [Charcoal], …)
//   - Reds & pinks ([Crimson], [Coral], [HotPink], …)
//   - Oranges & yellows ([Orange], [Gold], [Amber], …)
//   - Greens ([Lime], [Emerald], [Teal], …)
//   - Blues ([Navy], [RoyalBlue], [Cerulean], …)
//   - Purples & violets ([Violet], [Lavender], [Indigo], …)
//   - Browns & earth tones ([Brown], [Sienna], [Terracotta], …)
//   - Semantic UI aliases ([Success], [Warning], [Danger], [Info], [Muted])
var (
	// -----------------------------------------------------------------------
	// Basic colors
	// -----------------------------------------------------------------------

	// Black is pure black (0, 0, 0).
	Black = RGB(0, 0, 0)
	// White is pure white (255, 255, 255).
	White = RGB(255, 255, 255)
	// Red is a standard terminal red (205, 49, 49).
	Red = RGB(205, 49, 49)
	// Green is a standard terminal green (13, 188, 121).
	Green = RGB(13, 188, 121)
	// Yellow is a standard terminal yellow (229, 229, 16).
	Yellow = RGB(229, 229, 16)
	// Blue is a standard terminal blue (36, 114, 200).
	Blue = RGB(36, 114, 200)
	// Magenta is a standard terminal magenta (188, 63, 188).
	Magenta = RGB(188, 63, 188)
	// Cyan is a standard terminal cyan (17, 168, 205).
	Cyan = RGB(17, 168, 205)

	// -----------------------------------------------------------------------
	// Bright variants
	// -----------------------------------------------------------------------

	// BrightBlack is the bright variant of black, effectively a dark gray (102, 102, 102).
	BrightBlack = RGB(102, 102, 102)
	// BrightWhite is the bright variant of white (229, 229, 229).
	BrightWhite = RGB(229, 229, 229)
	// BrightRed is the bright variant of red (241, 76, 76).
	BrightRed = RGB(241, 76, 76)
	// BrightGreen is the bright variant of green (35, 209, 139).
	BrightGreen = RGB(35, 209, 139)
	// BrightYellow is the bright variant of yellow (245, 245, 67).
	BrightYellow = RGB(245, 245, 67)
	// BrightBlue is the bright variant of blue (59, 142, 234).
	BrightBlue = RGB(59, 142, 234)
	// BrightMagenta is the bright variant of magenta (214, 112, 214).
	BrightMagenta = RGB(214, 112, 214)
	// BrightCyan is the bright variant of cyan (41, 184, 219).
	BrightCyan = RGB(41, 184, 219)

	// -----------------------------------------------------------------------
	// Dark variants
	// -----------------------------------------------------------------------

	// DarkRed is a deep red (139, 0, 0).
	DarkRed = RGB(139, 0, 0)
	// DarkGreen is a deep green (0, 100, 0).
	DarkGreen = RGB(0, 100, 0)
	// DarkYellow is a muted dark yellow (155, 135, 12).
	DarkYellow = RGB(155, 135, 12)
	// DarkBlue is a deep blue (0, 0, 139).
	DarkBlue = RGB(0, 0, 139)
	// DarkMagenta is a deep magenta (139, 0, 139).
	DarkMagenta = RGB(139, 0, 139)
	// DarkCyan is a deep cyan (0, 139, 139).
	DarkCyan = RGB(0, 139, 139)
	// DarkGray is a dark neutral gray (64, 64, 64).
	DarkGray = RGB(64, 64, 64)
	// DarkOrange is a deep orange (255, 140, 0).
	DarkOrange = RGB(255, 140, 0)
	// DarkPurple is a deep indigo-purple (75, 0, 130).
	DarkPurple = RGB(75, 0, 130)
	// DarkPink is a deep rose pink (199, 21, 133).
	DarkPink = RGB(199, 21, 133)

	// -----------------------------------------------------------------------
	// Grays & neutrals
	// -----------------------------------------------------------------------

	// Gray is a mid-tone neutral gray (128, 128, 128).
	Gray = RGB(128, 128, 128)
	// Silver is a light metallic gray (192, 192, 192).
	Silver = RGB(192, 192, 192)
	// LightGray is a very light gray (211, 211, 211).
	LightGray = RGB(211, 211, 211)
	// DimGray is a slightly darker mid-gray (105, 105, 105).
	DimGray = RGB(105, 105, 105)
	// SlateGray is a blue-tinged gray (112, 128, 144).
	SlateGray = RGB(112, 128, 144)
	// Charcoal is a very dark blue-gray (54, 69, 79).
	Charcoal = RGB(54, 69, 79)
	// Gainsboro is a pale silver-gray (220, 220, 220).
	Gainsboro = RGB(220, 220, 220)
	// WhiteSmoke is an off-white near-white (245, 245, 245).
	WhiteSmoke = RGB(245, 245, 245)
	// OffWhite is a warm near-white with a faint cream tint (255, 250, 240).
	OffWhite = RGB(255, 250, 240)

	// -----------------------------------------------------------------------
	// Reds & pinks
	// -----------------------------------------------------------------------

	// Orange is a vivid orange (255, 165, 0).
	Orange = RGB(255, 165, 0)
	// Pink is a hot cerise pink (255, 105, 180).
	Pink = RGB(255, 105, 180)
	// Purple is a classic mid-purple (128, 0, 128).
	Purple = RGB(128, 0, 128)
	// Crimson is a rich deep red (220, 20, 60).
	Crimson = RGB(220, 20, 60)
	// Scarlet is a vivid orange-red (255, 36, 0).
	Scarlet = RGB(255, 36, 0)
	// Ruby is a deep gemstone red (155, 17, 30).
	Ruby = RGB(155, 17, 30)
	// Coral is a warm salmon-orange (255, 127, 80).
	Coral = RGB(255, 127, 80)
	// Salmon is a soft pink-orange (250, 128, 114).
	Salmon = RGB(250, 128, 114)
	// LightSalmon is a pale salmon (255, 160, 122).
	LightSalmon = RGB(255, 160, 122)
	// Tomato is a bright red-orange (255, 99, 71).
	Tomato = RGB(255, 99, 71)
	// HotPink is a vivid cerise pink (255, 105, 180).
	HotPink = RGB(255, 105, 180)
	// DeepPink is a strong magenta-pink (255, 20, 147).
	DeepPink = RGB(255, 20, 147)
	// LightPink is a pale pastel pink (255, 182, 193).
	LightPink = RGB(255, 182, 193)
	// RosyBrown is a muted dusty rose (188, 143, 143).
	RosyBrown = RGB(188, 143, 143)
	// Maroon is a dark brownish-red (128, 0, 0).
	Maroon = RGB(128, 0, 0)

	// -----------------------------------------------------------------------
	// Oranges & yellows
	// -----------------------------------------------------------------------

	// OrangeRed is a vivid red-leaning orange (255, 69, 0).
	OrangeRed = RGB(255, 69, 0)
	// Gold is a rich metallic yellow (255, 215, 0).
	Gold = RGB(255, 215, 0)
	// Amber is a warm yellow-orange (255, 191, 0).
	Amber = RGB(255, 191, 0)
	// Peach is a soft pale orange (255, 218, 185).
	Peach = RGB(255, 218, 185)
	// Khaki is a dull yellow-green (240, 230, 140).
	Khaki = RGB(240, 230, 140)
	// DarkKhaki is a muted olive-yellow (189, 183, 107).
	DarkKhaki = RGB(189, 183, 107)
	// LightYellow is a very pale yellow (255, 255, 224).
	LightYellow = RGB(255, 255, 224)
	// Lemon is a bright pure yellow (255, 247, 0).
	Lemon = RGB(255, 247, 0)
	// Wheat is a pale warm beige (245, 222, 179).
	Wheat = RGB(245, 222, 179)
	// Moccasin is a soft warm cream (255, 228, 181).
	Moccasin = RGB(255, 228, 181)

	// -----------------------------------------------------------------------
	// Greens
	// -----------------------------------------------------------------------

	// Lime is pure bright green (0, 255, 0).
	Lime = RGB(0, 255, 0)
	// LimeGreen is a vivid medium green (50, 205, 50).
	LimeGreen = RGB(50, 205, 50)
	// ForestGreen is a dark natural green (34, 139, 34).
	ForestGreen = RGB(34, 139, 34)
	// SeaGreen is a muted teal-green (46, 139, 87).
	SeaGreen = RGB(46, 139, 87)
	// MediumGreen is a solid mid-range green (0, 128, 0).
	MediumGreen = RGB(0, 128, 0)
	// SpringGreen is a vivid cyan-green (0, 255, 127).
	SpringGreen = RGB(0, 255, 127)
	// Chartreuse is a yellow-green (127, 255, 0).
	Chartreuse = RGB(127, 255, 0)
	// YellowGreen is a medium yellow-green (154, 205, 50).
	YellowGreen = RGB(154, 205, 50)
	// OliveGreen is a dark muted yellow-green (107, 142, 35).
	OliveGreen = RGB(107, 142, 35)
	// Olive is a dark yellow-green (128, 128, 0).
	Olive = RGB(128, 128, 0)
	// Teal is a dark cyan (0, 128, 128).
	Teal = RGB(0, 128, 128)
	// MintGreen is a pale bright green (152, 255, 152).
	MintGreen = RGB(152, 255, 152)
	// PaleGreen is a soft light green (152, 251, 152).
	PaleGreen = RGB(152, 251, 152)
	// DarkOlive is a very dark olive (85, 107, 47).
	DarkOlive = RGB(85, 107, 47)
	// Sage is a muted silvery green (143, 188, 143).
	Sage = RGB(143, 188, 143)
	// Emerald is a vivid jewel green (0, 201, 87).
	Emerald = RGB(0, 201, 87)
	// Jade is a deep cool green (0, 168, 107).
	Jade = RGB(0, 168, 107)

	// -----------------------------------------------------------------------
	// Blues
	// -----------------------------------------------------------------------

	// Navy is a very dark blue (0, 0, 128).
	Navy = RGB(0, 0, 128)
	// RoyalBlue is a rich medium blue (65, 105, 225).
	RoyalBlue = RGB(65, 105, 225)
	// SteelBlue is a muted steel blue (70, 130, 180).
	SteelBlue = RGB(70, 130, 180)
	// DodgerBlue is a vivid sky blue (30, 144, 255).
	DodgerBlue = RGB(30, 144, 255)
	// DeepSkyBlue is a bright clear sky blue (0, 191, 255).
	DeepSkyBlue = RGB(0, 191, 255)
	// SkyBlue is a soft light blue (135, 206, 235).
	SkyBlue = RGB(135, 206, 235)
	// LightBlue is a pale airy blue (173, 216, 230).
	LightBlue = RGB(173, 216, 230)
	// PowderBlue is a very light cool blue (176, 224, 230).
	PowderBlue = RGB(176, 224, 230)
	// CornflowerBlue is a medium periwinkle blue (100, 149, 237).
	CornflowerBlue = RGB(100, 149, 237)
	// MidnightBlue is a very dark navy blue (25, 25, 112).
	MidnightBlue = RGB(25, 25, 112)
	// Cerulean is a rich sky blue (0, 123, 167).
	Cerulean = RGB(0, 123, 167)
	// Cobalt is a deep saturated blue (0, 71, 171).
	Cobalt = RGB(0, 71, 171)
	// Sapphire is a deep gemstone blue (15, 82, 186).
	Sapphire = RGB(15, 82, 186)
	// Periwinkle is a soft blue-violet (204, 204, 255).
	Periwinkle = RGB(204, 204, 255)

	// -----------------------------------------------------------------------
	// Purples & violets
	// -----------------------------------------------------------------------

	// Violet is a pale purple (238, 130, 238).
	Violet = RGB(238, 130, 238)
	// Lavender is a very pale blue-purple (230, 230, 250).
	Lavender = RGB(230, 230, 250)
	// Orchid is a medium pink-purple (218, 112, 214).
	Orchid = RGB(218, 112, 214)
	// Plum is a soft muted purple (221, 160, 221).
	Plum = RGB(221, 160, 221)
	// Indigo is a deep blue-purple (75, 0, 130).
	Indigo = RGB(75, 0, 130)
	// SlateBlue is a medium blue-purple (106, 90, 205).
	SlateBlue = RGB(106, 90, 205)
	// MediumPurple is a soft medium purple (147, 112, 219).
	MediumPurple = RGB(147, 112, 219)
	// Amethyst is a rich gem purple (153, 102, 204).
	Amethyst = RGB(153, 102, 204)
	// Lilac is a pale pastel purple (200, 162, 200).
	Lilac = RGB(200, 162, 200)
	// Mauve is a pale violet (224, 176, 255).
	Mauve = RGB(224, 176, 255)
	// Fuchsia is a vivid magenta-pink (255, 0, 255).
	Fuchsia = RGB(255, 0, 255)
	// Thistle is a pale grey-purple (216, 191, 216).
	Thistle = RGB(216, 191, 216)

	// -----------------------------------------------------------------------
	// Browns & earth tones
	// -----------------------------------------------------------------------

	// Brown is a classic warm brown (165, 42, 42).
	Brown = RGB(165, 42, 42)
	// SaddleBrown is a dark reddish brown (139, 69, 19).
	SaddleBrown = RGB(139, 69, 19)
	// Sienna is a warm earthy brown (160, 82, 45).
	Sienna = RGB(160, 82, 45)
	// Peru is a medium earthy brown (205, 133, 63).
	Peru = RGB(205, 133, 63)
	// Chocolate is a deep warm brown (210, 105, 30).
	Chocolate = RGB(210, 105, 30)
	// Tan is a light warm brown (210, 180, 140).
	Tan = RGB(210, 180, 140)
	// BurlyWood is a pale wood brown (222, 184, 135).
	BurlyWood = RGB(222, 184, 135)
	// Sandy is a warm sandy orange (244, 164, 96).
	Sandy = RGB(244, 164, 96)
	// Bisque is a very pale warm beige (255, 228, 196).
	Bisque = RGB(255, 228, 196)
	// NavajoWhite is a soft warm cream (255, 222, 173).
	NavajoWhite = RGB(255, 222, 173)
	// Linen is a warm off-white (250, 240, 230).
	Linen = RGB(250, 240, 230)
	// Beige is a pale warm yellow-white (245, 245, 220).
	Beige = RGB(245, 245, 220)
	// Ivory is a near-white with a faint warm tint (255, 255, 240).
	Ivory = RGB(255, 255, 240)
	// Cornsilk is a very pale warm yellow (255, 248, 220).
	Cornsilk = RGB(255, 248, 220)
	// Umber is a dark earthy brown (99, 81, 71).
	Umber = RGB(99, 81, 71)
	// Taupe is a very dark warm gray-brown (72, 60, 50).
	Taupe = RGB(72, 60, 50)
	// Caramel is a warm golden brown (196, 127, 37).
	Caramel = RGB(196, 127, 37)
	// Rust is a deep red-orange brown (183, 65, 14).
	Rust = RGB(183, 65, 14)
	// Terracotta is a warm red clay (226, 114, 91).
	Terracotta = RGB(226, 114, 91)
	// Clay is a muted orange-red clay (211, 120, 89).
	Clay = RGB(211, 120, 89)

	// -----------------------------------------------------------------------
	// Special / UI-oriented
	// -----------------------------------------------------------------------

	// Success is a semantic alias for a bright green, suitable for success
	// messages and positive status indicators. Equivalent to [BrightGreen].
	Success = RGB(35, 209, 139)
	// Warning is a semantic alias for amber, suitable for warnings and caution
	// indicators. Equivalent to [Amber].
	Warning = RGB(255, 191, 0)
	// Danger is a semantic alias for bright red, suitable for errors and
	// critical status indicators. Equivalent to [BrightRed].
	Danger = RGB(241, 76, 76)
	// Info is a semantic alias for bright blue, suitable for informational
	// messages. Equivalent to [BrightBlue].
	Info = RGB(59, 142, 234)
	// Muted is a semantic alias for mid-gray, suitable for de-emphasized or
	// secondary text. Equivalent to [Gray].
	Muted = RGB(128, 128, 128)
)
