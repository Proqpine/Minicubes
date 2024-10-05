# frozen_string_literal: true

require_relative 'fnv1_has'
require 'digest'

# A Bloom filter is a space-efficient probabilistic data structure,
# conceived by Burton Howard Bloom in 1970, that is used to test whether
# an element is a member of a set. False positive matches are possible,
#  but false negatives are not – in other words, a query returns either
# "possibly in set" or "definitely not in set". Elements can be added to
# the set, but not removed (though this can be addressed with the counting
# Bloom filter variant); the more items added, the larger the probability
# of false positives.
# https://en.wikipedia.org/wiki/Bloom_filter
class BloomFilter
  attr_reader :data, :size_of_bits_array, :set_of_hash_functions, :num_of_entries

  def initialize(size_of_bits_array, set_of_hash_functions)
    @size_of_bits_array = size_of_bits_array
    @set_of_hash_functions = set_of_hash_functions
    @data = Array.new((size_of_bits_array + 7) / 8, 0)
    @num_of_entries = 0
  end

  def set_bit(position)
    byte_index = position / 8
    bit_index = position % 8
    @data[byte_index] |= (1 << bit_index)
  end

  def get_bit(position)
    byte_index = position / 8
    bit_index = position % 8
    (@data[byte_index] & (1 << bit_index)) != 0
  end

  # def generate_hashes(element)
  #   [
  #     Digest::SHA512.digest(element + '1').unpack1('Q*'),
  #     Digest::MD5.digest(element).unpack1('L*'),
  #     FNVHash.fnv_1a(element),
  #     Digest::SHA1.digest(element).unpack1('L*'),
  #     Digest::SHA256.digest(element).unpack1('L*'),
  #     Digest::SHA384.digest(element).unpack1('L*')
  #   ]
  # end

  # def insert(element)
  #   element = element.to_s.strip.downcase
  #   hash_values = generate_hashes(element)

  #   # Use as many hash values as needed based on @set_of_hash_functions
  #   hash_values.take(@set_of_hash_functions).each do |hash_val|
  #     set_bit(hash_val % @size_of_bits_array)
  #   end
  #   @num_of_entries += 1
  # end

  def insert(element)
    element = element.to_s.strip.downcase
    case @set_of_hash_functions
    when 1
      hash1 = Digest::SHA512.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      set_bit(hash1)
    when 2
      hash1 = Digest::SHA512.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      set_bit(hash1)
      set_bit(hash2)
    when 3
      hash1 = Digest::SHA512.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash3 = FNVHash.fnv_1a(element) % @size_of_bits_array
      set_bit(hash1)
      set_bit(hash2)
      set_bit(hash3)
    when 6
      hash1 = Digest::SHA512.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash3 = FNVHash.fnv_1a(element) % @size_of_bits_array
      hash4 = Digest::SHA1.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash5 = Digest::SHA256.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash6 = Digest::SHA384.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      set_bit(hash1)
      set_bit(hash2)
      set_bit(hash3)
      set_bit(hash4)
      set_bit(hash5)
      set_bit(hash6)
    end
    @num_of_entries += 1
  end

  def search(element)
    element = element.to_s.strip.downcase
    case @set_of_hash_functions
    when 1
      hash1 = Digest::SHA512.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      return "#{element} Not in Bloom Filter" unless get_bit(hash1)
    when 2
      hash1 = Digest::SHA512.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      return "#{element} Not in Bloom Filter" unless get_bit(hash1) && get_bit(hash2)
    when 3
      hash1 = Digest::SHA512.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash3 = FNVHash.fnv_1a(element) % @size_of_bits_array
      return "#{element} Not in Bloom Filter" unless get_bit(hash1) && get_bit(hash2) && get_bit(hash3)
    when 6
      hash1 = Digest::SHA512.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash3 = FNVHash.fnv_1a(element) % @size_of_bits_array
      hash4 = Digest::SHA1.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash5 = Digest::SHA256.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash6 = Digest::SHA384.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      return "#{element} Not in Bloom Filter" unless get_bit(hash1) && get_bit(hash2) && get_bit(hash3) &&
                                          get_bit(hash4) && get_bit(hash5) && get_bit(hash6)
    end

    prob = (1.0 - ((1.0 - (1.0 / @size_of_bits_array))**(@set_of_hash_functions * @num_of_entries)))**@set_of_hash_functions
    "#{element} Might be in Bloom Filter with false positive probability #{prob}"
  end
end
