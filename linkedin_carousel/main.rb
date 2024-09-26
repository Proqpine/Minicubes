# frozen_string_literal: true

require 'grover'

# LinkedIn Carousel Card
class Card
  NAME = 'Name'
  CONTENT = 'This is content'
  IMAGE = 'https://api.dicebear.com/9.x/lorelei/svg?seed=Brian'
  def initialize(name)
    @name = name
    @content = []
  end

  def read_template
    File.read('template.html')
  end

  def add_content=(content)
    @content << content
  end

  def add_profile_pic=(image)
    @image = image
  end

  def create_carousel(output_filename, content)
    template_content = read_template
    modified_content = template_content.gsub(/\b#{NAME}\b/, @name)
    modified_content = modified_content.gsub(/\b#{CONTENT}\b/, "\"#{content}\"")
    modified_content = modified_content.gsub(IMAGE, @image)
    if @template_content == modified_content
      puts 'Warning: No replacements were made. Check if placeholders exist in the template.'
    else
      puts 'Replacements made successfully.'
    end
    File.write(output_filename, modified_content)
    puts "Modified template has been saved as #{output_filename}"
    grover = Grover.new("file:///Users/davidayomide/Downloads/Dev/Minicubes/linkedin_carousel/#{output_filename}",
                        format: 'letter')

    image_name = File.basename(output_filename, '.html')
    File.write("#{image_name}.jpg", grover.to_jpeg)
  end

  def create_carousels
    num = 0
    @content.each do |content|
      create_carousel("#{num}.html", content)
      num += 1
    end
  end
end

# carousel = Card.new('Lame Ho')
# carousel.add_profile_pic = 'https://api.dicebear.com/9.x/notionists-neutral/svg?seed=Maria'
# carousel.add_content = 'We are the world'
# carousel.add_content = 'We of christ'

# carousel.create_carousels
grover = Grover.new('file:///Users/davidayomide/Downloads/Dev/Minicubes/linkedin_carousel/1.html',
                    format: 'A4')
File.open('groveer.jpg', 'wb') { |f| f << grover.to_jpeg }
