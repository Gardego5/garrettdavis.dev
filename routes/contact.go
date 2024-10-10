package routes

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"time"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/global"
	"github.com/Gardego5/garrettdavis.dev/middleware"
	"github.com/Gardego5/garrettdavis.dev/schema"
	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
	"github.com/go-playground/validator/v10"
)

func GetContact(w http.ResponseWriter, r *http.Request) {
	var name, email string
	switch rand.Intn(4) {
	case 0:
		name, email = "John Doe", "john.doe@mail.com"
	case 1:
		name, email = "Jane Doe", "jane.doe@gmail.com"
	case 2:
		name, email = "Mark Smith", "mark.smith@missivemark.dev"
	case 3:
		name, email = "Sara Clay", "sara@clay.dev"
	}

	middleware.RenderPage(r, Title{"Contact Garrett"},
		components.Header{},
		components.Margins(Form{Class("relative grid gap-2 rounded border border-slate-500 bg-gray-800 p-4 sm:grid-cols-2 md:grid-cols-3"),
			Attrs{{"hx-post", "/contact"}, {"hx-swap", "innerHTML"}, {"hx-target-error", "#form-error"}},
			H2{Class("text-2xl font-semibold col-span-full"),
				"I'd love to hear from you!",
			},

			Div{
				Label{Class("pl-2 text-slate-300 text-sm italic"), "Name"},
				Input{Class("block w-full px-2 py-1 rounded border bg-zinc-900 border-slate-500"),
					{"name", "name"}, {"required"}, {"type", "text"}, {"placeholder", name}},
				P{Class("text-sm text-red-600 px-2 [&:empty]:hidden"), Attr{"id", "name-error"}},
			},
			Div{Class("md:col-span-2"),
				Label{Class("pl-2 text-slate-300 text-sm italic"), "Email"},
				Input{Class("block w-full px-2 py-1 rounded border bg-zinc-900 border-slate-500 col-span-2"),
					{"name", "email"}, {"required"}, {"type", "email"}, {"placeholder", email}},
				P{Class("text-sm text-red-600 px-2 [&:empty]:hidden"), Attr{"id", "email-error"}},
			},
			Div{Class("col-span-full"),
				Label{Class("pl-2 text-slate-300 text-sm italic"), "Message"},
				Textarea{Attrs{Class("block w-full px-2 py-1 rounded border bg-zinc-900 border-slate-500"),
					{"name", "message"}, {"required"}, {"placeholder", "¡Hej! ¿What do you think about this?"}}},
				P{Class("text-sm text-red-600 px-2 [&:empty]:hidden"), Attr{"id", "message-error"}},
			},

			Input{{"type", "hidden"}, {"name", "dummy-name"}, {"value", name}},
			Input{{"type", "hidden"}, {"name", "dummy-email"}, {"value", email}},

			P{Class("col-span-full text-sm text-red-600 px-2 [&:empty]:hidden"), Attr{"id", "form-error"}},

			Button{Attrs{Class("absolute -bottom-[7px] right-8 rounded border border-slate-500 bg-zinc-900 px-4 py-1 text-sm hover:bg-slate-800"),
				{"type", "submit"}},
				"Send",
			},
		}),
	)
}

type contactMessageErrors struct {
	General, Name, Email, Message error
}

func (c contactMessageErrors) Render() RenderedHTML {
	return Fragment{
		Span{Switch().
			Case(c.General != nil, Block(func() any { return c.General.Error() })).
			Case(c.Name == nil && c.Email == nil && c.Message == nil, "Some unkown error occurred. Please try again.").
			Default("Please fix these errors.")},
		Div{Attr{"hx-swap-oob", "innerHTML:#name-error"},
			Span{If(c.Name != nil, Block(func() any { return c.Name.Error() }))}},
		Div{Attr{"hx-swap-oob", "innerHTML:#email-error"},
			Span{If(c.Email != nil, Block(func() any { return c.Email.Error() }))}},
		Div{Attr{"hx-swap-oob", "innerHTML:#message-error"},
			Span{If(c.Message != nil, Block(func() any { return c.Message.Error() }))}},
	}
}

func PostContact(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	db := middleware.GetDB(r)

	body := schema.ContactMessage{
		Name:      r.FormValue("name"),
		Email:     r.FormValue("email"),
		Message:   r.FormValue("message"),
		CreatedAt: time.Now(),
	}

	if err := global.Validate.Struct(body); err != nil {
		resp := contactMessageErrors{}

		for _, e := range err.(validator.ValidationErrors) {
			fieldName := e.StructField()
			err := fmt.Errorf("%s is required.", fieldName)
			reflect.ValueOf(&resp).
				Elem().
				FieldByName(fieldName).
				Set(reflect.ValueOf(err))
		}

		logger.Warn("validation error", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		Render(w, resp)
		return
	}

	if _, err := db.NamedExec("INSERT INTO contact_messages (name, email, message, created_at) VALUES (:name, :email, :message, :created_at)", body); err != nil {
		logger.Error("error inserting contact message", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		Render(w, contactMessageErrors{General: err})
		return
	}

	Render(w, Div{Class("col-span-full text-center text-xl"),
		"Thanks ", body.Name, ", I'll get back to you soon."})
}
