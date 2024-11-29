# garrettdavis.dev

My personal site, and a demostration of building an htmx based app using
aws lambda, api gateway and cloudfront.

## TODO

### Infrastructure

- [ ] Cachix
- [ ] S3 Terraform State
- [ ] API Gateway
- [ ] Cloudfront
- [ ] Static Files
- [ ] SyncThing
      To get SyncThing working, we might need an ec2 instance that will host a
      SyncThing server to recieve all the data.

### Pages

- [x] Landing Page
- [ ] Notes Browser
- [ ] Projects Page
- [ ] About Me
- [x] Resume Page
- [ ] Blog Index
      somewhat done, ready, but only displayed on index currently
- [ ] Coffee Today

## Coffee Today

A short form more regular way to post something.
Easier talking point is what coffee did I have today, sort of a mini review...
Probably won't be daily, but a good way to keep the site fresh.

### Logged In

- [ ] Ability to create coffee posts with simple form.
- [ ] Prompt user for images to upload.
- [ ] Preview before posting.
- [ ] Automatically add the date.

Schema akin to:

```go
type Coffee struct {
    Date time.Time
    Title *string // optional - only replace default title if present
    Roaster string
    RoasterUrl *string
    TastingNotes []string // optional
    Pictures []string // store in tigris s3

    Ideas []string // optional - short notes or ideas for the day, maybe like a 'did you know?'
}
```

### Public

- [ ] Users can view an index of coffee posts by date.
- [ ] The index shows the images from the post, title, and date.
- [ ] Homepage has a coffee mug that points to the most recent coffee post.
- [ ] Send me a coffee button that links to some way to send funds.
- [ ] Contact form for sending me a message...
      Maybe different from normal contact form, and prompt for optionally including a picture of your coffee.
