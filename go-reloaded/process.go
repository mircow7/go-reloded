package main

import (
	"strconv"
	"strings"
	"unicode"
)

// ProcessText corrige tout le texte, ligne par ligne
// (comme ça les retours à la ligne ne disparaissent pas)
func ProcessText(text string) string {
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = Process(lines[i])
	}
	return strings.Join(lines, "\n")
}

// Process corrige une seule ligne, en 4 étapes (l'ordre est important !)
func Process(line string) string {
	// Étape 1 : les marqueurs (hex), (bin), (up), (low), (cap)
	words := splitLine(line)
	words = applyMarkers(words)

	// Étape 2 : la ponctuation
	text := fixPunctuation(strings.Join(words, " "))

	// Étape 3 : les guillemets
	words = fixQuotes(strings.Fields(text))

	// Étape 4 : a -> an
	words = fixArticles(words)

	return strings.Join(words, " ")
}

// ---------------------------------------------------------------
// ÉTAPE 1 : les marqueurs
// ---------------------------------------------------------------

// isNumber dit si le texte est fait uniquement de chiffres
func isNumber(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// readMarker regarde ce qu'il y a entre les parenthèses.
// "cap, 6" -> ("cap", 6, true)   "up" -> ("up", 1, true)   "hello" -> faux
func readMarker(inside string) (string, int, bool) {
	inside = strings.ReplaceAll(inside, " ", "")
	parts := strings.Split(inside, ",")

	cmd := parts[0]
	if cmd != "hex" && cmd != "bin" && cmd != "up" && cmd != "low" && cmd != "cap" {
		return "", 0, false
	}

	// pas de nombre : on prend 1 par défaut
	if len(parts) == 1 {
		return cmd, 1, true
	}
	if len(parts) > 2 || !isNumber(parts[1]) {
		return "", 0, false
	}

	n, err := strconv.Atoi(parts[1])
	if err != nil {
		n = 1000000 // nombre énorme : on modifiera tous les mots
	}
	return cmd, n, true
}

// tokenMarker teste si un mot est un marqueur comme "(up,2)"
func tokenMarker(tok string) (string, int, bool) {
	if len(tok) < 3 || tok[0] != '(' || tok[len(tok)-1] != ')' {
		return "", 0, false
	}
	return readMarker(tok[1 : len(tok)-1])
}

// splitLine découpe la ligne en mots.
// Particularité : un marqueur comme "(cap, 6)" contient un espace,
// donc on le garde en un seul morceau, écrit proprement "(cap,6)".
func splitLine(line string) []string {
	var tokens []string
	current := ""

	i := 0
	for i < len(line) {
		c := line[i]

		// Est-ce le début d'un marqueur ?
		if c == '(' {
			end := strings.Index(line[i:], ")")
			if end != -1 {
				inside := line[i+1 : i+end]
				cmd, n, ok := readMarker(inside)
				if ok {
					if current != "" {
						tokens = append(tokens, current)
						current = ""
					}
					tokens = append(tokens, "("+cmd+","+strconv.Itoa(n)+")")
					i = i + end + 1
					continue
				}
			}
		}

		// Sinon : espace = fin du mot, autre caractère = on l'ajoute au mot
		if c == ' ' || c == '\t' || c == '\r' {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		} else {
			current = current + line[i:i+1]
		}
		i++
	}

	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens
}

// isWord : un vrai mot, pas juste de la ponctuation ou un guillemet
func isWord(tok string) bool {
	for _, c := range tok {
		if !strings.ContainsRune(".,!?:;'", c) {
			return true
		}
	}
	return false
}

// capitalize met la 1re lettre en majuscule et le reste en minuscules
func capitalize(w string) string {
	letters := []rune(strings.ToLower(w))
	if len(letters) == 0 {
		return w
	}
	letters[0] = unicode.ToUpper(letters[0])
	return string(letters)
}

// changeWord applique une commande sur un mot
func changeWord(cmd string, w string) string {
	if cmd == "up" {
		return strings.ToUpper(w)
	}
	if cmd == "low" {
		return strings.ToLower(w)
	}
	if cmd == "cap" {
		return capitalize(w)
	}
	if cmd == "hex" {
		n, err := strconv.ParseInt(w, 16, 64)
		if err != nil {
			return w // pas un nombre valide : on laisse le mot
		}
		return strconv.FormatInt(n, 10)
	}
	if cmd == "bin" {
		n, err := strconv.ParseInt(w, 2, 64)
		if err != nil {
			return w
		}
		return strconv.FormatInt(n, 10)
	}
	return w
}

// applyMarkers parcourt les mots. Quand on tombe sur un marqueur,
// on revient en arrière pour modifier les n mots d'avant,
// et le marqueur n'est pas gardé.
func applyMarkers(tokens []string) []string {
	var result []string

	for _, tok := range tokens {
		cmd, n, ok := tokenMarker(tok)
		if !ok {
			result = append(result, tok)
			continue
		}

		// on remonte dans result, on saute la ponctuation
		for i := len(result) - 1; i >= 0 && n > 0; i-- {
			if !isWord(result[i]) {
				continue
			}
			result[i] = changeWord(cmd, result[i])
			n--
		}
	}
	return result
}

// ---------------------------------------------------------------
// ÉTAPE 2 : la ponctuation
// ---------------------------------------------------------------

func isPunct(r rune) bool {
	return strings.ContainsRune(".,!?:;", r)
}

// fixPunctuation : pas d'espace avant . , ! ? : ;  et une espace après.
// Les groupes comme "..." ou "!?" restent ensemble.
func fixPunctuation(text string) string {
	runes := []rune(text)
	var result []rune

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// espace juste avant une ponctuation : on l'enlève
		if r == ' ' && i+1 < len(runes) && isPunct(runes[i+1]) {
			continue
		}

		result = append(result, r)

		// ponctuation collée à une lettre : on ajoute une espace après
		if isPunct(r) && i+1 < len(runes) {
			next := runes[i+1]
			if next != ' ' && !isPunct(next) {
				result = append(result, ' ')
			}
		}
	}
	return string(result)
}

// ---------------------------------------------------------------
// ÉTAPE 3 : les guillemets
// ---------------------------------------------------------------

// fixQuotes : ' ouvrant collé au mot d'après, ' fermant collé au mot d'avant.
// On ne regarde que les ' tout seuls, donc don't n'est pas touché.
func fixQuotes(words []string) []string {
	var result []string
	open := false    // est-ce qu'on est dans une citation ?
	waiting := false // un ' ouvrant attend le prochain mot

	for _, w := range words {
		if w == "'" {
			if !open {
				open = true
				waiting = true
			} else {
				open = false
				if len(result) > 0 {
					result[len(result)-1] = result[len(result)-1] + "'"
				} else {
					result = append(result, "'")
				}
			}
			continue
		}

		if waiting {
			w = "'" + w
			waiting = false
		}
		result = append(result, w)
	}

	if waiting {
		result = append(result, "'")
	}
	return result
}

// ---------------------------------------------------------------
// ÉTAPE 4 : a -> an
// ---------------------------------------------------------------

// fixArticles : "a" devient "an" si le mot d'après commence par
// a, e, i, o, u ou h  (pareil pour A -> An)
func fixArticles(words []string) []string {
	for i := 0; i < len(words)-1; i++ {
		if words[i] != "a" && words[i] != "A" {
			continue
		}

		next := strings.TrimLeft(words[i+1], "'") // on ignore un ' ouvrant
		if next == "" {
			continue
		}

		first := unicode.ToLower([]rune(next)[0])
		if strings.ContainsRune("aeiouh", first) {
			words[i] = words[i] + "n"
		}
	}
	return words
}
