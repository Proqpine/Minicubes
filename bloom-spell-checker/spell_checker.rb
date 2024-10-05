# frozen_string_literal: true

require 'optparse'
require_relative 'helper'
require_relative 'bloom_filter'

class BloomFilerSpellChecker < BloomFilter
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

  def initialize
    super(0, 0)
  end

  def build_bloom_filter(file_name)
    helper = Helper.new
    size = helper.determine_size(file_name)
    num_of_hash = helper.determine_number_of_hash

    @size_of_bits_array = size
    @set_of_hash_functions = num_of_hash
    @data = Array.new((size + 7) / 8, 0)
    @num_of_entries = 0

    head_chuck = Header.new('SCBF', 1, num_of_hash, size)

    puts "Building Bloom filter with size: #{size}, hash functions: #{num_of_hash}"

    words_processed = 0
    BloomFilter.new(size, num_of_hash)
    File.open('words.bf', 'wb') do |file|
      file.write([head_chuck.id].pack('A4'))
      file.write([head_chuck.version_num].pack('n'))
      file.write([head_chuck.num_hash_function].pack('n'))
      file.write([head_chuck.bloom_filter_size].pack('N'))

      File.foreach(file_name) do |line|
        word = line.strip
        next if word.empty?

        insert(word)
        words_processed += 1
        puts "Processed #{words_processed} words..." if words_processed % 1000 == 0
      end

      file.write(@data.pack('C*'))
    end

    puts "Completed building Bloom filter with #{words_processed} words"
  rescue StandardError => e
    puts "Error building Bloom filter: #{e.message}"
    puts e.backtrace
    raise
  end

  def read_bloom_filter(file_name)
    File.open(file_name, 'rb') do |file|
      # Read header information
      id = file.read(4)
      version = file.read(2).unpack1('n')
      num_hash = file.read(2).unpack1('n')
      bloom_filter_size = file.read(4).unpack1('N')

      puts "ID: #{id}, Version: #{version}, Hash Functions: #{num_hash}, Size: #{bloom_filter_size}"

      # Initialize the bloom filter with the correct parameters
      @set_of_hash_functions = num_hash
      @size_of_bits_array = bloom_filter_size
      @data = Array.new(@size_of_bits_array, 0)

      # Read the bloom filter data
      bloom_data = file.read
      @data = bloom_data.unpack('C*')
      @num_of_entries = bloom_filter_size

      # If there are command line arguments, search for them
      if ARGV.any?
        ARGV.each do |word|
          result = search(word.strip)
          puts "#{word}: #{result}"
        end
      end
    end
  rescue StandardError => e
    puts "Error reading Bloom filter: #{e.message}"
    puts e.backtrace
    raise
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
