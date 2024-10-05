# frozen_string_literal: true

expected_elements = 260_000  # Number of expected unique elements
false_positive_rate = 0.01   # Desired false positive rate

# Calculate size of bits array
size_of_bits_array = (-expected_elements * Math.log(false_positive_rate)) / (Math.log(2)**2)
size_of_bits_array = size_of_bits_array.ceil # Round up

# Convert bits to bytes
byte_size = (size_of_bits_array + 7) / 8 # Add 7 to round up to the nearest byte

puts "Expected Elements: #{expected_elements}"
puts "Calculated Size of Bits Array: #{size_of_bits_array} bits"
puts "Calculated Size in Bytes: #{byte_size} bytes"
