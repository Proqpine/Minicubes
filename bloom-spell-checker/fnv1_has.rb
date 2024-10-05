# frozen_string_literal: true

# The FNV-1a hash differs from the FNV-1 hash only by the order in which the multiply and XOR is performed
# https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function
class FNVHash
  FNV_PRIME_32 = 16_777_619
  FNV_OFFSET_32 = 2_166_136_261
  def fnv_1a(input_data)
    bytes = input_data.is_a?(String) ? input_data.bytes : input_data
    hash = FNV_OFFSET_32
    bytes.each do |byte|
      hash ^= byte
      hash *= FNV_PRIME_32
      hash &= 0xffffffff
    end
    hash
  end
end

puts fnv_1a('hello world')
