# WordList Ruby Class

A simple Ruby class for loading and managing the Wordle word list.

## Usage

```ruby
require_relative 'word_list'

# Create a new word list instance
word_list = WordList.new

# Get basic information
puts "Total words: #{word_list.size}"
puts "Random word: #{word_list.random_word}"

# Check if a word is valid
if word_list.include?('hello')
  puts "'hello' is a valid word"
end

# Get words of specific length
five_letter_words = word_list.five_letter_words
puts "Number of five-letter words: #{five_letter_words.size}"

# Get words of any length
three_letter_words = word_list.words_of_length(3)
puts "Number of three-letter words: #{three_letter_words.size}"
```

## Methods

- `initialize(file_path = nil)` - Creates a new WordList instance
- `size` - Returns the total number of words
- `include?(word)` - Checks if a word is in the list (case-insensitive)
- `random_word` - Returns a random word from the list
- `words_of_length(length)` - Returns all words of the specified length
- `five_letter_words` - Returns all five-letter words
- `reload!` - Reloads the word list from the file
- `to_a` - Returns the words as an array
- `to_set` - Returns the words as a Set

## Running the Test

To see the WordList in action, run:

```bash
ruby test_word_list.rb
```

This will show:
- Total word count
- Random word samples
- Five-letter word count and samples
- Word validation examples
- Word length distribution
- Sample words by length

## File Structure

- `valid-wordle-words.txt` - The source word list file
- `word_list.rb` - The WordList class implementation
- `test_word_list.rb` - Test script demonstrating usage
- `README.md` - This documentation file
