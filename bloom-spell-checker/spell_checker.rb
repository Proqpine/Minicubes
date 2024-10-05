# frozen_string_literal: true

require 'optparse'
require_relative 'helper'
require_relative 'bloom_filter'

class BloomFilerSpellChecker
  options = {}
  Header = Struct.new(:id, :version_num, :num_hash_function, :bloom_filter_size)

  OptionParser.new do |opts|
    opts.banner = 'Usage: spell_checker.rb [options]'

    opts.on('-b', '--build FILE', 'Create a .bf bloom filter') do |bf|
      options[:build] = bf
    end
  end.parse!

  def build_bloom_filter(file_name)
    helper = Helper.new
    size = helper.determine_size(file_name)
    _num_of_hash = helper.determine_number_of_hash

    head_chuck = Header.new('SCBF', 1, 3, size)

    bloom_filter = BloomFilter.new(size, 3)
    File.open('words.bf', 'wb') do |file|
      file.write([head_chuck.id].pack('A4'))
      file.write([head_chuck.version_num].pack('n'))
      file.write([head_chuck.num_hash_function].pack('n'))
      file.write([head_chuck.bloom_filter_size].pack('N'))

      File.foreach(file_name).with_index do |element, index|
        bloom_filter.insert(element.strip)
        puts "Inserted #{index + 1} entries..." if (index + 1) % 1000 == 0
      end
      file.write(bloom_filter.data.pack('C*'))
    end
  end

  if options[:build]
    bf = BloomFilerSpellChecker.new
    bf.build_bloom_filter(options[:build])
  else
    puts 'No file specified. Use -build FILE to create a file.'
  end
end
