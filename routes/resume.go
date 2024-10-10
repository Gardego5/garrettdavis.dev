package routes

import (
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/middleware"
	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
	"github.com/elliotchance/pie/v2"
)

type (
	descBlock struct{ string }
	descList  []string
	item      struct {
		annotation, title, subtitle string
		desc                        HTML
	}
	section struct {
		title string
		items []item
	}
)

func (d descBlock) Render() RenderedHTML {
	return P{Class("text-sm"), d.string}.Render()
}

func (d descList) Render() RenderedHTML {
	return Ul{Class("flex gap-x-3 gap-y-2 flex-wrap text-sm"),
		pie.Map(d, func(s string) HTML { return Li{s} }),
	}.Render()
}

func (i item) Render() RenderedHTML {
	return Div{Class("pb-3"),
		H3{Class("text-blue-300 text-xl"),
			i.title,

			If(i.annotation != "", Span{Class("inline-block float-right"),
				i.annotation,
			}),
		},

		If(i.subtitle != "", H4{Class("text-slate-300 text-sm pl-2 -mt-1 mb-1 italic print:text-xs"),
			i.subtitle,
		}),

		i.desc,
	}.Render()
}

func (s section) Render() RenderedHTML {
	f := Fragment{
		H2{Class("font-mono leading-7 tracking-lighter"), s.title},
		Div{pie.Map(s.items, func(i item) any { return i })},
	}
	return f
}

type resumePair [2]any

func (r resumePair) Render() RenderedHTML {
	return Fragment{
		H2{Class("font-mono leading-7 tracking-lighter col-start-1"), r[0]},
		Div{Class("col-start-2"), r[1]},
	}
}

func GetResume(w http.ResponseWriter, r *http.Request) {
	middleware.RenderPage(r,
		Fragment{
			Title{"Resume - Garrett Davis"},
			Meta{{"name", "description"}, {"content", "Garrett Davis' resume"}},
		},

		Div{Class("fixed top-0 right-0 flex gap-2 p-2 print:hidden"),
			Button{Attrs{{"x-data"}, {"x-on:click", "window.print()"}, {"title", "Print this page."}},
				Element("iconify-icon", Attrs{{"icon", "ph:printer"}, {"width", "36"}, {"height", "36"}}),
			},
		},

		Div{Class("mx-12"),
			Div{Class("grid grid-cols-[1fr_6fr] [&>*:nth-child(odd)]:text-right gap-4 print:text-black max-w-4xl m-auto"),
				Div{Class("col-start-2 flex py-4 items-center"),
					H1{
						Span{Class("text-4xl font-mono font-bold"), "Garrett Davis"},
						Br{},
						Span{Class("text-lg")},
					},
				},

				P{Class("col-start-2"),
					"I am a software engineer. I love to make reliable, maintainable, simple, and enjoyable to use software.",
				},

				resumePair{
					"Contact",
					Ul{
						Li{"Email: ", A{Attrs{
							{"x-data", `{user:'contact',domain:window.location.hostname}`},
							{"x-init", `$el.href = 'mailto:' + $data.user + '@' + $data.domain`},
						},
							"contact [at] ", Span{Attrs{{"x-text", "domain"}}},
						}},

						Li{"Location: ", "Hillsboro, Oregon"},
					},
				},

				section{title: "Technical Skills", items: []item{
					{title: "Programming Languages", desc: descList{
						"Go",
						"TypeScript",
						"JavaScript",
						"Nix",
						"Rust",
						"HCL",
						"Python",
						"Shell",
						"PHP",
						"Java",
					}},
					{title: "Frontend", desc: descList{
						"NextJS",
						"React",
						"Vue",
						"Nuxt",
						"HTML",
						"CSS",
						"TailwindCSS",
						"Bootstrap",
						"MaterialUI",
						"SASS",
					}},
					{title: "Libraries", desc: descList{
						"Sqlx",
						"Gin",
						"Go's testing package",
						"Express",
						"NextAuth",
						"TRPC",
						"Axum",
						"Zod",
						"esbuild",
					}},
					{title: "Databases", desc: descList{
						"MySQL",
						"DynamoDB",
						"SQLite",
						"PostgreSQL",
						"IndexedDB",
					}},
					{title: "CI/CD Tools", desc: descList{
						"Terraform",
						"AWS CodeBuild",
						"AWS CodePipeline",
						"GitHub Actions",
						"NixOS",
					}},
				}},

				section{title: "Personal Skills", items: []item{
					{
						title: "Teaching",
						desc:  descBlock{"I seek out new technologies and look for how they might benefit my current work. After I find something useful and applicable, I love to share this with my team - and strive to help my whole group excel. For example while working as a contractor at University of Phoenix, I volunteered and then lead a talk on building AWS Lambda functions using Rust."},
					},
					{
						title: "Leadership",
						desc:  descBlock{"I managed a team of 3-7 baristas while working at Starbucks. I lead our store to be the top in our district for customer connection by focusing on technical excellence and meaningful conversations with our customers."},
					},
				}},

				section{title: "Experience", items: []item{
					{
						title:      "University of Phoenix",
						subtitle:   "Software Engineer I",
						annotation: "10/23 - Present",
						desc: descList{
							"I spearheaded the transition from typescript to Go for backend microservices.",
							"I built a custom role based authorization system supporting multiple organizations and scopes.",
							"When requirements changed, I refactored our Go backend from microservices to a monolith architecture to support the new use case.",
							"I built multiple smooth login processes using OAuth2 integrating with AzureAD, AWS Cognito, and Google auth providers.",
						},
					},
					{
						title:      "Cook Systems - University of Phoenix",
						subtitle:   "Contract Software Engineer",
						annotation: "10/22 - 10/23",
						desc: descList{
							"Revamped existing authentication flow built with PHP to integrate with NextAuth and custom Django / EdX authentication system.",
							"I built automatically generated sitemap for marketing website.",
							"I built a tool to programatically test lighthouse scores across entire marketing site.",
							"I built a custom replacement for Terraform Cloud using s3 backend and codepipeline, reducing time to create a new project's cicd infrastrucure from roughly 8 hours down to only 15 - 30 minutes.",
						},
					},
					{
						title:      "MSR-FSR",
						subtitle:   "Production Technician",
						annotation: "10/21 - 02/22",
						desc: descList{
							"Performed detail oriented work in a cleanroom environment.",
						},
					},
					{
						title:      "Starbucks",
						subtitle:   "Shift Supervisor",
						annotation: "08/18 - 10/21",
						desc: descList{
							"Mentored multiple baristas toward promotion to supervisor by encouraging them to coach others.",
							"I created opportunities for training and used commendation to encourage individual growth.",
							"Grew our  team by training 18+ new hires while instilling Starbucks quality and customer service values.",
							"I lead coffee tastings focusing on the science and history in each cup.",
							"Managed a team of 3-7 people, promoting communication and teamwork.",
						},
					},
				}},

				section{title: "Education", items: []item{
					{
						title:      "Cook Systems FastTrack'D",
						annotation: "07/22 - 09/22",
						desc:       descBlock{"Concentrated Java Frameworks and developer tools training."},
					},
					{
						title:      "Portland Community College",
						annotation: "09/16 - 03/20",
						desc:       descBlock{"Associate of General Studies. Emphasis on GIS, Cartography, and Math."},
					},
				}},
			},
		},
	)
}
