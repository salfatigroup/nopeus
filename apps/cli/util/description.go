package util

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	colors "gopkg.in/go-playground/colors.v1"
)

// generate a list of gradient points between two colors
func getGradientList(startColor string, endColor string, steps uint) ([]string, error) {
    start, err := colors.ParseHEX(startColor)
    if err != nil {
        return nil, err
    }

    end, err := colors.ParseHEX(endColor)
    if err != nil {
        return nil, err
    }

    var colorsList []string
    rdelta := float64(end.ToRGB().R - start.ToRGB().R) / float64(steps)
    gdelta := float64(end.ToRGB().G - start.ToRGB().G) / float64(steps)
    bdelta := float64(end.ToRGB().B - start.ToRGB().B) / float64(steps)

    for i := uint(0); i < steps; i++ {
        r := start.ToRGB().R + uint8(rdelta*float64(i))
        g := start.ToRGB().G + uint8(gdelta*float64(i))
        b := start.ToRGB().B + uint8(bdelta*float64(i))
        newColor := fmt.Sprintf("#%02x%02x%02x", r, g, b)
        colorsList = append(colorsList, newColor)
    }

    return colorsList, nil
}

// make the given text in gradient colors
func GradientText(text string, startColor string, endColor string) string {
    var result string
    colorsList, err := getGradientList(startColor, endColor, uint(len(text)))
    if err != nil {
        return text
    }

    for i, c := range text {
        result += lipgloss.NewStyle().
            Bold(true).
            Foreground(lipgloss.Color(colorsList[i])).
            Render(string(c))
    }

    return result
}

// Convert text to a gray color
func GrayText(text string) string {
    return lipgloss.NewStyle().Faint(true).Render(text)
}

// Return a colorful short description to the terminal
// cli commander
func ShortDescription() string {
    return lipgloss.NewStyle().
        MarginTop(1).
        Render(GradientText("Nopeus - Cloud Application Layer", "#db2777", "#f9a8d4"))
}

// Return a long description of nopeus to the terminal
func LongDescription() string {
    return ShortDescription() + "\n\n" + lipgloss.NewStyle().Faint(true).Render("Nopeus provides an application layer to your cloud infrastructure. Shift left infrastructure provisioning. Nopeus provides an opiniated tool that aims to simplify cloud provisioning by providing an application layer to the cloud. Nopeus's goal it to ensure a scalable and secure infrastructure with minimum configurations (and will always try to remove more configurations then addition). Nopeus is designed with monorepo and microservices in mind, but can work in any structure.")
}
