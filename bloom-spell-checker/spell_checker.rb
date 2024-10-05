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
    opts.on('-r', '--read FILE', 'Read a .bf bloom filter') do |rbf|
      options[:read] = rbf
    end
  end.parse!

  def build_bloom_filter(file_name)
    helper = Helper.new
    size = helper.determine_size(file_name)
    num_of_hash = helper.determine_number_of_hash

    head_chuck = Header.new('SCBF', 1, num_of_hash, size)

    bloom_filter = BloomFilter.new(size, num_of_hash)
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

  def read_bloom_filter(file_name)
    File.open(file_name, 'rb') do |file|
      id = file.read(4)
      version = file.read(2).unpack1('n')
      num_hash = file.read(2).unpack1('n')
      bloom_filter_size = file.read(4).unpack1('N')

      puts "ID: #{id}, Version: #{version}, Hash Functions: #{num_hash}, Size: #{bloom_filter_size}"

      # Read the rest of the file (the Bloom filter data)
      file.read(bloom_filter_size)
      # Additional logic to handle the bloom_filter_data
    end
  end

  if options[:build]
    bf = BloomFilerSpellChecker.new
    bf.build_bloom_filter(options[:build])
  elsif options[:read]
    bf = BloomFilerSpellChecker.new
    bf.read_bloom_filter(options[:read])
  else
    puts 'No file specified. Use -build FILE to create a file.'
  end
end
