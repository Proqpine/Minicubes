# frozen_string_literal: true

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

def encode_alphanumeric_data(input)
  data_array = input.chars.map { |char| ALPHANUMERIC_CHARS.find_index(char) }
  encoded_bits = []

  data_array.each_slice(2) do |pair|
    if pair.length == 2
      pair_value = pair[0] * 45 + pair[1]
      encoded_bits << pair_value.to_s(2).rjust(11, '0')
    else
      encoded_bits << pair[0].to_s(2).rjust(6, '0')
    end
  end

  encoded_bits.join
end

def encode_numeric_data(input)
  encoded_bits = []
  input.chars.each_slice(3) do |three|
    if three.length == 3
      value = three.join.to_i
      encoded_bits << value.to_s(2).rjust(10, '0')
    elsif three.length == 2
      value = three.join.to_i
      encoded_bits << value.to_s(2).rjust(7, '0')
    elsif three.length == 1
      value = three.join.to_i
      encoded_bits << value.to_s(2).rjust(4, '0')
    end
  end
  encoded_bits.join
end

def encode_byte_data(input)
  encoded_bits = []
  data_array = input.chars.map(&:ord)
  data_array.each do |dar|
    encoded_bits << dar.to_s(2).rjust(8, '0')
  end
  encoded_bits.join
end

def determine_input_mode(input)
  if numeric?(input)
    mode = :numeric
    char_count_indicator = encode_character_count(input, mode)
    puts "Input: #{input}"
    puts 'Mode: Numeric'
    puts "Character Count: #{input.length}"
    puts "Character Count Indicator (binary): #{char_count_indicator}"
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
