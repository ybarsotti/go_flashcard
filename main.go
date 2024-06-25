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

type Card struct {
	Term       string
	Definition string
	answered	bool
}

func NewCard(term, definition string) *Card {
	return &Card{
		Term:       term,
		Definition: definition,
	}
}

func readInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return input[:len(input)-1], nil
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

func promptAction() string {
	var action string
	fmt.Println("Input the action (add, remove, import, export, ask, exit):")
	fmt.Scan(&action)
	return action
}

func createFlashCard(cards *[]Card) {
	card := Card{}
	fmt.Println("The card:")
	for {
		term, err := readInput(bufio.NewReader(os.Stdin))

		if err != nil {
			fmt.Errorf("Error parsing card. %w", err)
			continue
		}

		if is_created, _ := isCardTermAlreadyCreated(*cards, term); is_created {
			fmt.Println("This term already exists. Try again:")
			continue
		}

		card.Term = term
		break
	}

	fmt.Println("The definition of the card:")
	for {
		definition, err := readInput(bufio.NewReader(os.Stdin))
		if err != nil {
			fmt.Errorf("Error parsing definition. %w", err)
			continue
		}

		if isCardDefinitionAlreadyCreated(*cards, definition) {
			fmt.Println("This definition already exists. Try again:")
			continue
		}
		card.Definition = definition
		break
	}
	*cards = append(*cards, card)
	fmt.Printf("The pair (\"%s\":\"%s\") has been added\n", card.Term, card.Definition)
}

func removeFlashCard(cards *[]Card) {
	fmt.Println("Which card?")
	term, err := readInput(bufio.NewReader(os.Stdin))
	if err != nil {
		fmt.Errorf("Error reading term. %w", err)
	}

	for index, card := range *cards {
		if card.Term == term {
			*cards = append((*cards)[:index], (*cards)[index+1:]...)
			fmt.Println("The card has been removed.")
			return
		}
	}
	fmt.Printf("Can't remove \"%s\": there is no such card.\n", term)
}

func exportFlashCard(cards *[]Card) {
	fmt.Println("File name:")
	filename, err := readInput(bufio.NewReader(os.Stdin))
	if err != nil {
		fmt.Errorf("Error reading filename. %w", err)
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Errorf("Error reading file. %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "term,definition\n")

	for _, card := range *cards {
		_, err = fmt.Fprintf(file, "%s,%s\n", card.Term, card.Definition)
		if err != nil {
			fmt.Errorf("Error writing line to file %w", err)
		}
	}
	fmt.Printf("%d cards have been saved.\n", len(*cards))
}

func importFlashCard(cards *[]Card) {
	fmt.Println("File name:")
	filename, err := readInput(bufio.NewReader(os.Stdin))
	if err != nil {
		fmt.Printf("Error reading filename. %s", err)
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("File not found.\n")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	var loaded_count int
	for scanner.Scan() {
		card_data := strings.Split(scanner.Text(), ",")
		term, definition := card_data[0], card_data[1]
		if is_created, index := isCardTermAlreadyCreated(*cards, term); is_created {
			// Delete card
			*cards = append((*cards)[:index], (*cards)[index+1:]...)
		}
		*cards = append(*cards, *NewCard(term, definition))
		loaded_count += 1
	}
	fmt.Printf("%d cards have been loaded.\n", loaded_count)
}

func askCards(cards *[]Card) {
	fmt.Println("How many times to ask?")
	amount, err := readInput(bufio.NewReader(os.Stdin))
	if err != nil {
		fmt.Printf("Error parsing count %s", err)
	}
	
	amount_int, err := strconv.Atoi(amount)
	if err != nil {
		fmt.Printf("Error parsing count %s", err)
	}

	if len(*cards) < 1 {
		fmt.Println("Not enough cards.")
		return
	}
	
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < amount_int; i++ {
		card := (*cards)[r.Intn(len(*cards) - 1)]
		fmt.Printf("Print the definition of \"%s\":\n", card.Term)
		answer, err := readInput(bufio.NewReader(os.Stdin))
		if err != nil {
			fmt.Printf("Error parsing answer %s", err)
			continue
		}

		if card.Definition == answer {
			fmt.Println("Correct!")
		} else {
			var is_correct_for_other bool
			fmt.Print("Wrong. ")
			for _, other_card := range *cards {
				if other_card.Definition == answer {
					fmt.Printf("The right answer is \"%s\", but your definition is correct for \"%s\".\n", card.Definition, other_card.Term)
					is_correct_for_other = true
					break
				}
			}

			if !is_correct_for_other {
				fmt.Printf("The right answer is \"%s\".\n", card.Definition)
			}
		}
	}

}

func main() {
	cards := []Card{}

promptFor:
	for {
		action := promptAction()
		switch action {
		case "add":
			createFlashCard(&cards)
		case "remove":
			removeFlashCard(&cards)
		case "import":
			importFlashCard(&cards)
		case "export":
			exportFlashCard(&cards)
		case "ask":
			askCards(&cards)
		case "exit":
			fmt.Println("Bye bye!")
			break promptFor
		}
	}

}
