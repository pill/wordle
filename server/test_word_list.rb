#!/usr/bin/env ruby

require_relative 'word_list'

# Create a new word list instance
word_list = WordList.new

puts "=== WordList Test ==="
puts "Total words loaded: #{word_list.size}"
puts "Random word: #{word_list.random_word}"
puts "Random word: #{word_list.random_word}"
puts "Random word: #{word_list.random_word}"

puts "\n=== Five Letter Words ==="
five_letter_words = word_list.five_letter_words
puts "Number of five-letter words: #{five_letter_words.size}"
puts "Sample five-letter words: #{five_letter_words.sample(10).join(', ')}"

puts "\n=== Word Validation ==="
test_words = ['hello', 'world', 'apple', 'xyzzy', 'abask', 'aahed']
test_words.each do |word|
  valid = word_list.include?(word)
  puts "'#{word}' is #{valid ? 'valid' : 'not valid'}"
end

puts "\n=== Word Length Distribution ==="
(3..8).each do |length|
  count = word_list.words_of_length(length).size
  puts "#{length}-letter words: #{count}"
end

puts "\n=== Sample Words by Length ==="
(3..8).each do |length|
  words = word_list.words_of_length(length)
  if words.any?
    sample = words.sample(5).join(', ')
    puts "#{length}-letter words sample: #{sample}"
  end
end 