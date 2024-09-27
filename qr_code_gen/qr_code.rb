# frozen_string_literal: false

data_array = []
'The quick brown fox'.split('').map { |char| data_array.push(char.ord) }
puts data_array.length

# = Modes
# Numeric
# Alphanumeric
# Kanji

NUMERIC = 'Numeric'
ALPHANUMERIC = 'Alphanumeric'
KANJI = 'Kanji'
BYTE = 'Byte'
ALPHANUMERIC_CHARS = ('0'..'9').to_a + ('A'..'Z').to_a + [' ', '$', '%', '*', '+', '-', '.', '/', ':']

def numeric?(input)
  input.match?(/\A[0-9]+\z/)
end

def alphanumeric?(input)
  input.chars.all? { |char| ALPHANUMERIC_CHARS.include?(char) }
end

def kanji?(input)
  # This is a simplified check. A more accurate check would involve
  # checking against a comprehensive list of Kanji characters or using a gem.
  input.encode('UTF-8').match?(/\p{Han}|\p{Hiragana}|\p{Katakana}/)
rescue Encoding::UndefinedConversionError
  false
end

def determine_input_mode(input)
  if numeric?(input)
    NUMERIC
  elsif alphanumeric?(input)
    ALPHANUMERIC
  elsif kanji?(input)
    KANJI
  else
    BYTE
  end
end

# Test the function
test_inputs = [
  '12345',
  'HELLO WORLD',
  'こんにちは',
  'Hello, World!',
  '42',
  'QR CODE',
  'ｱｲｳｴｵ', # Katakana
  'Mixed 123 Content!'
]

test_inputs.each do |input|
  puts "Input: #{input}"
  puts "Mode: #{determine_input_mode(input)}"
  puts '---'
end
