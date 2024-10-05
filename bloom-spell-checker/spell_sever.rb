# frozen_string_literal: true

require 'sinatra'
require 'json'
require_relative 'spell_checker'

get '/' do
  erb :index
end

post '/submit' do
  data = JSON.parse(request.body.read)
  word = data['word']

  bf = BloomFilerSpellChecker.new
  bf.read_bloom_filter('words.bf')
  out = bf.search(word)

  if word.empty?
    'Please enter your name.'
  else
    "#{out}!"
  end
end

# server = WEBrick::HTTP::Server.new(Port: 8000)

# server.mount_proc '/' do |_req, res|
#   res.body = '<html><body><h1>Hello, World!</h1></body></html>'
# end
