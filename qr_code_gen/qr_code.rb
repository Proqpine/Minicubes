# frozen_string_literal: true

require 'reedsolomon'

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

MODE_INDICATOR = {
  numeric: 0b0001,
  alphanumeric: 0b0010,
  byte: 0b0100,
  kanji: 0b1000
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

def add_mode_bits(arr, mode, input)
  arr << MODE_INDICATOR[mode].to_s(2).rjust(4, '0') # Mode Indicator
  arr << encode_character_count(input, mode) # Character count
end

def byte_align(cou_bits)
  unless (cou_bits.length % 8).zero?
    padding_needed = 8 - (cou_bits.length % 8)
    cou_bits += '0' * padding_needed
  end
  cou_bits
end

def add_padding(ec_level, cou_bits)
  first_pad = 0xEC.to_s(2).rjust(8, '0') # Ensure 8 bits
  second_pad = 0x11.to_s(2).rjust(8, '0')

  # Define the size based on the error correction level
  size = case ec_level
         when 'L' then 80 * 8
         when 'M' then 64 * 8
         when 'Q' then 48 * 8
         when 'H' then 36 * 8
         else 80 * 8
         end

  # Calculate how much padding is needed
  padding_left = size - cou_bits.length

  # Alternate padding bytes 0xEC and 0x11
  alternate = padding_left / 8 # Divide by 8 since we add bytes

  cou_bits_arr = []
  cou_bits_arr << cou_bits

  alternate.times do |i|
    cou_bits_arr << if i.even?
                      first_pad # Add 0xEC
                    else
                      second_pad # Add 0x11
                    end
  end
  cou_bits_arr.join
end

def encode_full_string(input)
  full_bits = []
  if numeric?(input)
    mode = :numeric
    puts NUMERIC
    add_mode_bits(full_bits, mode, input)
    full_bits << encode_numeric_data(input) # Data
  elsif alphanumeric?(input)
    mode = :alphanumeric
    puts ALPHANUMERIC
    add_mode_bits(full_bits, mode, input)
    full_bits << encode_alphanumeric_data(input)
  else
    mode = :byte
    puts BYTE
    add_mode_bits(full_bits, mode, input)
    full_bits << encode_byte_data(input)
  end
  full_bits << '0000' # Add terminator
  full_bits.join
  cou_bits = full_bits.join.to_s
  byte_align(cou_bits)
  add_padding('L', cou_bits)
end

puts encode_full_string('HELLO WORLD')
