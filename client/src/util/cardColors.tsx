enum CardColor {
    RED = "rgba(230, 35, 25, 1)",
    YELLOW = "rgba(245, 245, 10, 1)",
    GREEN = "rgba(40, 210, 25, 1)",
    LIGHT_BLUE = "rgba(120, 245, 250, 1)",
    DARK_BLUE = "rgba(30, 65, 255, 1)"
}

export const cardColorLookup: {[key: number]: string} = {
    [-2]: CardColor.DARK_BLUE,
    [-1]: CardColor.DARK_BLUE,
    [0]: CardColor.LIGHT_BLUE,
    [1]: CardColor.GREEN,
    [2]: CardColor.GREEN,
    [3]: CardColor.GREEN,
    [4]: CardColor.GREEN,
    [5]: CardColor.YELLOW,
    [6]: CardColor.YELLOW,
    [7]: CardColor.YELLOW,
    [8]: CardColor.YELLOW,
    [9]: CardColor.RED,
    [10]: CardColor.RED,
    [11]: CardColor.RED,
    [12]: CardColor.RED,
};