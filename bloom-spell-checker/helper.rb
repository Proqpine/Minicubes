# frozen_string_literal: true

# Determines the number of items that are likely to be stored and
# the probability of false positives the system can tolerate then
# uses that to determine the memory requirements and number of hash
# functions needed.
class Helper
  attr_reader :num_of_entries, :array_size

  FALSE_POSITIVE_RATE = 0.01 # Set a default value for false positive rate

  def initialize
    @num_of_entries = 0
    @array_size = 0
  end

  def determine_size(input_data)
    raise ArgumentError, 'File does not exist' unless File.exist?(input_data)

    file = File.open(input_data, 'r')
    @num_of_entries = file.readlines.size
    @array_size = calculate_array_size
  end

  def determine_number_of_hash
    (Math.log(2) * (@array_size / @num_of_entries)).round
  end

  def calculate_array_size
    return 0 if @num_of_entries.zero?

    -(@num_of_entries * Math.log(FALSE_POSITIVE_RATE) / (Math.log(2)**2)).round
  end
end

helper = Helper.new
puts helper.determine_size('dict.txt')
puts helper.determine_number_of_hash
