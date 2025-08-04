class WordList
    attr_reader :words
  
    def initialize(file_path = nil)
      @file_path = file_path || File.join(__dir__, 'valid-wordle-words.txt')
      @words = load_words
    end
  
    def load_words
      return [] unless File.exist?(@file_path)
      
      File.readlines(@file_path).map(&:strip).reject(&:empty?)
    end
  
    def size
      @words.size
    end
  
    def include?(word)
      @words.include?(word.downcase)
    end
  
    def random_word
      @words.sample
    end
  
    def words_of_length(length)
      @words.select { |word| word.length == length }
    end
  
    def five_letter_words
      words_of_length(5)
    end
  
    def reload!
      @words = load_words
    end
  
    def to_a
      @words.dup
    end
  
    def to_set
      Set.new(@words)
    end
  end
  
  # Example usage:
  # word_list = WordList.new
  # puts "Total words: #{word_list.size}"
  # puts "Random word: #{word_list.random_word}"
  # puts "Five letter words: #{word_list.five_letter_words.size}"
  # puts "Is 'hello' valid? #{word_list.include?('hello')}"