package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

type Card struct {
	Term       string
	Definition string
	Mistakes   int
}

func (c *Card) IncrementMistake() {
	c.Mistakes += 1
}

func (c *Card) ResetMistakes() {
	c.Mistakes = 0
}

type LoggerAcc struct {
	b      strings.Builder
	reader *bufio.Reader
}

func NewLoggerAcc(builder *strings.Builder, reader *bufio.Reader) *LoggerAcc {
	return &LoggerAcc{
		b:      *builder,
		reader: reader,
	}
}

func (lc *LoggerAcc) ReadInput() (string, error) {
	input, err := lc.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	lc.b.Write([]byte("> " + input))
	return input[:len(input)-1], nil
}

func (lc *LoggerAcc) addOutput(output string) {
	lc.b.Write([]byte(output))
}

func (lc *LoggerAcc) OutputToConsole(text string, a ...any) (int, error) {
	parsedText := fmt.Sprintf(text+"\n", a...)
	lc.addOutput(parsedText)
	return fmt.Printf(parsedText)
}

func (lc *LoggerAcc) OutputToConsoleInline(text string, a ...any) (int, error) {
	parsedText := fmt.Sprintf(text, a...)
	lc.addOutput(parsedText)
	return fmt.Printf(parsedText)
}

func (lc *LoggerAcc) OutputData() string {
	return lc.b.String()
}

func NewCard(term, definition string, mistakes int) *Card {
	return &Card{
		Term:       term,
		Definition: definition,
		Mistakes:   mistakes,
	}
}

func isCardTermAlreadyCreated(cards []Card, term string) (bool, int) {
	for index, card := range cards {
		if card.Term == term {
			return true, index
		}
	}
	return false, -1
}

func isCardDefinitionAlreadyCreated(cards []Card, definition string) bool {
	for _, card := range cards {
		if card.Definition == definition {
			return true
		}
	}
	return false
}

func promptAction(loggerAcc *LoggerAcc) string {
	text := "Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):"
	loggerAcc.OutputToConsole(text)
	action, _ := loggerAcc.ReadInput()
	return action
}

func createFlashCard(cards *[]Card, loggerAcc *LoggerAcc) {
	card := Card{}
	loggerAcc.OutputToConsole("The card:")
	for {
		term, err := loggerAcc.ReadInput()

		if err != nil {
			loggerAcc.OutputToConsole("Error parsing card. %s", err)
			continue
		}

		if is_created, _ := isCardTermAlreadyCreated(*cards, term); is_created {
			loggerAcc.OutputToConsole(Red + "This term already exists. Try again:" + Reset)
			continue
		}

		card.Term = term
		break
	}

	loggerAcc.OutputToConsole("\nThe definition of the card:")
	for {
		definition, err := loggerAcc.ReadInput()
		if err != nil {
			loggerAcc.OutputToConsole("Error parsing definition. %s", err)
			continue
		}

		if isCardDefinitionAlreadyCreated(*cards, definition) {
			loggerAcc.OutputToConsole("This definition already exists. Try again:")
			continue
		}
		card.Definition = definition
		break
	}
	*cards = append(*cards, card)
	loggerAcc.OutputToConsole(Green+"The pair (\"%s\":\"%s\") has been added\n"+Reset, card.Term, card.Definition)
}

func removeFlashCard(cards *[]Card, loggerAcc *LoggerAcc) {
	loggerAcc.OutputToConsole("Which card?")
	term, err := loggerAcc.ReadInput()
	if err != nil {
		loggerAcc.OutputToConsole("Error reading term. %s", err)
	}

	for index, card := range *cards {
		if card.Term == term {
			*cards = append((*cards)[:index], (*cards)[index+1:]...)
			loggerAcc.OutputToConsole(Green + "The card has been removed.\n" + Reset)
			return
		}
	}
	loggerAcc.OutputToConsole("Can't remove \"%s\": there is no such card.\n", term)
}

func exportFlashCard(cards *[]Card, loggerAcc *LoggerAcc) {
	loggerAcc.OutputToConsole("File name:")
	filename, err := loggerAcc.ReadInput()
	if err != nil {
		loggerAcc.OutputToConsole("Error reading filename. %s", err)
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		loggerAcc.OutputToConsole("Error reading file. %s", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "term,definition,mistakes\n")

	for _, card := range *cards {
		_, err = fmt.Fprintf(file, "%s,%s,%d\n", card.Term, card.Definition, card.Mistakes)
		if err != nil {
			loggerAcc.OutputToConsole("Error writing line to file %s", err)
		}
	}
	loggerAcc.OutputToConsole(Green+"%d cards have been saved.\n"+Reset, len(*cards))
}

func importFlashCard(cards *[]Card, loggerAcc *LoggerAcc) {
	loggerAcc.OutputToConsole("File name:")
	filename, err := loggerAcc.ReadInput()
	if err != nil {
		loggerAcc.OutputToConsole("Error reading filename. %s", err)
	}
	file, err := os.Open(filename)
	if err != nil {
		loggerAcc.OutputToConsole("File not found.")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	var loaded_count int
	for scanner.Scan() {
		card_data := strings.Split(scanner.Text(), ",")
		term, definition, mistakes := card_data[0], card_data[1], card_data[2]
		if is_created, index := isCardTermAlreadyCreated(*cards, term); is_created {
			// Delete card
			*cards = append((*cards)[:index], (*cards)[index+1:]...)
		}
		mistakes_int, _ := strconv.Atoi(mistakes) // TODO: Check this?
		*cards = append(*cards, *NewCard(term, definition, mistakes_int))
		loaded_count += 1
	}
	loggerAcc.OutputToConsole(Green+"%d cards have been loaded.\n"+Reset, loaded_count)
}


func askCards(cards *[]Card, loggerAcc *LoggerAcc) {
	loggerAcc.OutputToConsole("How many times to ask?")
	amount, err := loggerAcc.ReadInput()
	if err != nil {
		loggerAcc.OutputToConsole("Error parsing count %s", err)
	}

	amount_int, err := strconv.Atoi(amount)
	if err != nil {
		loggerAcc.OutputToConsole("Error parsing count %s", err)
	}

	if len(*cards) < 1 {
		loggerAcc.OutputToConsole("Not enough cards.")
		return
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < amount_int; i++ {
		card_position := r.Intn(len(*cards)-1)
		card := (*cards)[card_position]
		loggerAcc.OutputToConsole("Print the definition of \"%s\":", card.Term)
		answer, err := loggerAcc.ReadInput()
		if err != nil {
			loggerAcc.OutputToConsole("Error parsing answer %s", err)
			continue
		}

		if card.Definition == answer {
			loggerAcc.OutputToConsole(Green + "Correct!\n" + Reset)
		} else {
			(*cards)[card_position].IncrementMistake()
			var is_correct_for_other bool
			for _, other_card := range *cards {
				if other_card.Definition == answer {
					loggerAcc.OutputToConsole(Magenta+"Wrong. The right answer is \"%s\", but your definition is correct for \"%s\".\n"+Reset, card.Definition, other_card.Term)
					is_correct_for_other = true
					break
				}
			}

			if !is_correct_for_other {
				loggerAcc.OutputToConsole(Red+"Wrong. The right answer is \"%s\".\n"+Reset, card.Definition)
			}
		}
	}

}

func saveLogs(loggerAcc *LoggerAcc) {
	loggerAcc.OutputToConsole("File name:")
	filename, err := loggerAcc.ReadInput()
	if err != nil {
		loggerAcc.OutputToConsole("Error reading filename. %s", err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		loggerAcc.OutputToConsole("Error reading file. %s", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, loggerAcc.OutputData())
	if err != nil {
		loggerAcc.OutputToConsole("Error writing line to file %s", err)
		return
	}

	loggerAcc.OutputToConsole(Green + "The log has been saved." + Reset)
}

func startLoggerBuilder() *strings.Builder {
	var b strings.Builder
	return &b
}

func listHardestCards(cards *[]Card, loggerAcc *LoggerAcc) {
	var higherMistakes int
	var termsWithMostMistakes []string

	for _, card := range *cards {
		if card.Mistakes > higherMistakes {
			higherMistakes = card.Mistakes
		}
	}

	for _, card := range *cards {
		if card.Mistakes == higherMistakes {
			termsWithMostMistakes = append(termsWithMostMistakes, card.Term)
		}
	}

	termsCount := len(termsWithMostMistakes)
	if higherMistakes == 0 {
		loggerAcc.OutputToConsole("There are no cards with errors.")
		return
	} 
	if termsCount == 1 {
		loggerAcc.OutputToConsole("The hardest card is \"%s\". You have %d errors answering it", termsWithMostMistakes[0], higherMistakes)
		return
	} 

	loggerAcc.OutputToConsoleInline("The hardest cards are ")
	
	for index, term := range termsWithMostMistakes {
		if index == 0 {
			loggerAcc.OutputToConsoleInline("\"%s\"", term)
		} else {
			loggerAcc.OutputToConsoleInline(", \"%s\"", term)
		}
		
	}
	loggerAcc.OutputToConsoleInline("\n")
}

func resetStats(cards *[]Card, loggerAcc *LoggerAcc) {
	for index, _ := range *cards {
		(*cards)[index].ResetMistakes()
	}
	loggerAcc.OutputToConsole("Card statistics have been reset.")
}

func main() {
	cards := []Card{}
	loggerAcc := NewLoggerAcc(startLoggerBuilder(), bufio.NewReader(os.Stdin))

promptFor:
	for {
		action := promptAction(loggerAcc)
		switch action {
		case "add":
			createFlashCard(&cards, loggerAcc)
		case "remove":
			removeFlashCard(&cards, loggerAcc)
		case "import":
			importFlashCard(&cards, loggerAcc)
		case "export":
			exportFlashCard(&cards, loggerAcc)
		case "ask":
			askCards(&cards, loggerAcc)
		case "log":
			saveLogs(loggerAcc)
		case "hardest card": 
			listHardestCards(&cards, loggerAcc)
		case "reset stats":
			resetStats(&cards, loggerAcc)
		case "exit":
			loggerAcc.OutputToConsole(Blue + "Bye bye!" + Reset)
			break promptFor
		}
	}

}
