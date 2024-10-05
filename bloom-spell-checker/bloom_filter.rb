# frozen_string_literal: true

require_relative 'fnv1_has'
require 'digest'

FALSE_POSITIVE_RATE = 0.01

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
  def initialize(size_of_bits_array, set_of_hash_functions)
    @size_of_bits_array = size_of_bits_array
    @set_of_hash_functions = set_of_hash_functions
    @data = [0] * @size_of_bits_array
    @num_of_entries = 0
  end

  attr_reader :data

  def insert(element)
    if @set_of_hash_functions == 1
      hash1 = element.hash.abs % @size_of_bits_array
      @data[hash1] = 1
    elsif @set_of_hash_functions == 2
      hash1 = element.hash.abs % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*')[0] % @size_of_bits_array
      @data[hash1] = 1
      @data[hash2] = 1
    elsif @set_of_hash_functions == 3
      hash1 = element.hash.abs % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*')[0] % @size_of_bits_array
      hash3 = FNVHash.fnv_1a(element) % @size_of_bits_array
      @data[hash1] = 1
      @data[hash2] = 1
      @data[hash3] = 1
    else
      raise "Unsupported number of hash functions: #{@set_of_hash_functions}"
    end
    @num_of_entries += 1
  end

  def search(element)
    if @set_of_hash_functions == 1
      hash1 = element.hash.abs % @size_of_bits_array
      return 'Not in Bloom Filter' if @data[hash1] == 0
    elsif @set_of_hash_functions == 2
      hash1 = element.hash.abs % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      return 'Not in Bloom Filter' if @data[hash1] == 0 || @data[hash2] == 0
    elsif @set_of_hash_functions == 3
      hash1 = element.hash.abs % @size_of_bits_array
      hash2 = Digest::MD5.digest(element.to_s).unpack1('L*') % @size_of_bits_array
      hash3 = FNVHash.fnv_1a(element) % @size_of_bits_array
      return 'Not in Bloom Filter' if @data[hash1] == 0 || @data[hash2] == 0 || @data[hash3] == 0
    end
    prob = (1.0 - ((1.0 - (1.0 / @size_of_bits_array))**(@set_of_hash_functions * @num_of_entries)))**@set_of_hash_functions
    # prob = (1 - Math.exp(-@set_of_hash_functions * @num_of_entries.to_f / @size_of_bits_array))**@set_of_hash_functions

    "Might be in Bloom Filter with false positive probability #{prob}"
  end
end
