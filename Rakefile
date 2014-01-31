S3Bucket = "workmachine.us"

desc "Deploy the frontend"
task :deploy_frontend do
  Dir.chdir "frontend" do
    # Run Jekyll
    puts "-> Running Jekyll"
    system "jekyll build"

    puts "\n\n-> Uploading to S3"

    # Sync media files first (Cache: expire in 10weeks)
    puts "\n--> Syncing media files..."
    system %{s3cmd sync --acl-public --exclude '*.*' --include '*.png' --include '*.jpg' --include '*.ico' --add-header="Expires: Sat, 20 Nov 2020 18:46:39 GMT" --add-header="Cache-Control: max-age=6048000"  _site/ s3://#{S3Bucket}/}

    # Sync Javascript and CSS assets next (Cache: expire in 1 week)
    puts "\n--> Syncing .js and .css files..."
    system %{s3cmd sync --acl-public --exclude '*.*' --include  '*.css' --include '*.js' --add-header="Cache-Control: max-age=3600"  _site/ s3://#{S3Bucket}}

    # Sync html files (Cache: 2 hours)
    puts "\n--> Syncing .html"
    system %{s3cmd sync --acl-public --exclude '*.*' --include  '*.html' --add-header="Cache-Control: max-age=3600, must-revalidate"  _site/ s3://#{S3Bucket}}

    # Sync everything else, but ignore the assets!
    puts "\n--> Syncing everything else"
    system %{s3cmd sync --acl-public --exclude '.DS_Store' --exclude '/personal/'  _site/ s3://#{S3Bucket}/}

    # Sync: remaining files & delete removed
    system %{s3cmd sync --acl-public --delete-removed  _site/ s3://#{S3Bucket}/}
  end
end

task :deploy => [:deploy_frontend] do

end
