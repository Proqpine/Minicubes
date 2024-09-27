# frozen_string_literal: true

# Modes
# Numeric
# Alphanumeric
# Kanji

NUMERIC = 'Numeric'
ALPHANUMERIC = 'Alphanumeric'
KANJI = 'Kanji'
BYTE = 'Byte'
ALPHANUMERIC_CHARS = ('0'..'9').to_a + ('A'..'Z').to_a + [' ', '$', '%', '*', '+', '-', '.', '/', ':']

CHARACTER_COUNT_BITS = {
  numeric: 10,
  alphanumeric: 9,
  byte: 8,
  kanji: 8
}

def numeric?(input)
  input.match?(/\A[0-9]+\z/)
end

def alphanumeric?(input)
  input.chars.all? { |char| ALPHANUMERIC_CHARS.include?(char) }
end

def kanji?(input)
  input.encode('UTF-8').match?(/\p{Han}|\p{Hiragana}|\p{Katakana}/)
rescue Encoding::UndefinedConversionError
  false
end

def encode_character_count(input, mode)
  char_count = input.length
  bit_count = CHARACTER_COUNT_BITS[mode]

  char_count.to_s(2).rjust(bit_count, '0')
end

def determine_input_mode(input)
  if numeric?(input)
    NUMERIC
  elsif alphanumeric?(input)
    mode = :alphanumeric
    char_count_indicator = encode_character_count(input, mode)
    puts "Input: #{input}"
    puts 'Mode: Alphanumeric'
    puts "Character Count: #{input.length}"
    puts "Character Count Indicator (binary): #{char_count_indicator}"
  elsif kanji?(input)
    KANJI
  else
    BYTE
  end
end

# EC_LEVEL 15%
# Version 4 Size 33 * 33
# Mode Indicator Alphanumeric 0010
# Character count binary

input_string = 'HELLO CC WORLD'

data_array = input_string.chars.map { |char| ALPHANUMERIC_CHARS.find_index(char) }
encoded_bits = []

data_array.each_slice(2) do |pair|
  if pair.length == 2
    puts "Pair: #{pair}"
    pair_value = pair[0] * 45 + pair[1]
    encoded_bits << pair_value.to_s(2).rjust(11, '0')
  else
    encoded_bits << pair[0].to_s(2).rjust(6, '0')
  end
end

string_out = encoded_bits.join
puts "Data: #{string_out}"

# 0110000101101111000110100010111000100010100011001110100100010100110111011111
# 01100001011011110001101000101110001000101000110011101001000101001101110111110
