# frozen_string_literal: true

require 'grover'

# LinkedIn Carousel Card
class Card
  NAME = 'Name'
  CONTENT = 'This is content'
  IMAGE = 'https://api.dicebear.com/9.x/lorelei/svg?seed=Brian'
  def initialize(name, num_pages)
    @name = name
    @num_pages = num_pages
  end

  def read_template
    File.read('template.html')
  end

  def add_content=(content)
    @content = content
  end

  def add_profile_pic=(image)
    @image = image
  end

  def create_carousel(output_filename)
    template_content = read_template
    modified_content = template_content.gsub(/\b#{NAME}\b/, @name)
    modified_content = modified_content.gsub(/\b#{CONTENT}\b/, "\"#{@content}\"")
    modified_content = modified_content.gsub(IMAGE, @image)
    if @template_content == modified_content
      puts 'Warning: No replacements were made. Check if placeholders exist in the template.'
    else
      puts 'Replacements made successfully.'
    end
    File.write(output_filename, modified_content)
    puts "Modified template has been saved as #{output_filename}"
  end
end

carousel = Card.new('Lame Ho', 1)
carousel.add_content = 'We are the world'
carousel.add_profile_pic = 'https://api.dicebear.com/9.x/notionists-neutral/svg?seed=Maria'
carousel.create_carousel('a.html')

# Grover.new accepts a URL or inline HTML and optional parameters for Puppeteer
grover = Grover.new('http://127.0.0.1:5500/template.html', format: 'letter')

# Get a screenshot
File.write('out.jpg', grover.to_jpeg)
