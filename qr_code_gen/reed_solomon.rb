# frozen_string_literal: true

class ReedSolomon
  FIELD_SIZE = 256
  PRIMITIVE_POLYNOMIAL = 285

  def add_in_field(num_a, num_b)
    num_a ^ num_b
  end

  def multiply_in_field(num_a, num_b)
    result = 0
    while num_b > 0
      result ^= num_a if (num_b & 1) != 0
      num_a <<= 1
      num_a ^= 0x11b if num_a & 0x100 != 0
      num_b >>= 1
    end
    result
  end

  def polynomial_multiply(p1, p2)
    result = [0] * (p1.length + p2.length - 1)
    (0...p1.length - 1).each do |i|
      (0...p2.length - 1).each do |j|
        result[i + j] = add_in_field(result[i + j], multiply_in_field(p1[i], p2[j]))
      end
    end
    result
  end

  def exp_in_field(alpha, power)
    result = 1
    power.times do
      result = multiply_in_field(result, alpha)
    end
    result
  end

  def generate_generator_polynomial(num_of_ec_cw)
    g = [1] # Start with x^0 = 1
    (0...num_of_ec_cw).each do |i|
      # Create (x + α^i)
      alpha_i = exp_in_field(2, i) # α^i in GF(2^8)
      factor = [alpha_i, 1] # α^i * x^0 + 1 * x^1
      g = polynomial_multiply(g, factor)
    end
    g
  end

  def encode(message, n_check_symbols)
    # Input validation
    raise ArgumentError, 'Message cannot be empty' if message.empty?
    raise ArgumentError, 'Number of check symbols must be positive' if n_check_symbols <= 0

    if message.length + n_check_symbols > FIELD_SIZE
      raise ArgumentError,
            'Message length + check symbols must not exceed field size'
    end

    generator = generate_generator_polynomial(n_check_symbols)
    encoded = message + [0] * n_check_symbols
    (0...message.length - 1).each do |i|
      coeff = encoded[i]
      next unless coeff != 0

      (0...generator.length - 1).each do |j|
        encoded[i + j] = add_in_field(encoded[i + j], multiply_in_field(generator[j], coeff))
      end
    end
    encoded
  end
end
