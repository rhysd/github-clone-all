guard :shell do
  watch /\.go$/ do |m|
    puts "#{Time.now}: #{m[0]}"
    case m[0]
    when /_test\.go$/
      system 'go test'
    else
      system 'go build'
    end
  end
end
