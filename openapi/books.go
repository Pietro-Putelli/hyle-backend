package openapi

import (
	"net/http"

	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
)

// BuildBooksAPI builds the books openapi API endpoints
func BuildBooksAPI(reflector openapi3.Reflector, securityName string) error {
	// BooksGetFun
	booksGetOp, err := reflector.NewOperationContext(http.MethodGet, "/v1/books/{bookID}")
	if err != nil {
		return err
	}
	booksGetOp.AddSecurity(securityName)
	booksGetOp.SetSummary("Retrieve book details by UUID")
	booksGetOp.SetDescription("Retrieve details of a book by its UUID")
	booksGetOp.SetID("bookDetailsGet")
	// booksGetOp.AddReqStructure(domain.BookGetPath{})
	booksGetOp.AddRespStructure(domain.BookResponse{})
	booksGetOp.AddRespStructure(failure.Error{}, openapi.WithHTTPStatus(http.StatusBadRequest))
	booksGetOp.SetTags("Books")

	err = reflector.AddOperation(booksGetOp)
	if err != nil {
		return err
	}

	// BooksListFun
	booksListOp, err := reflector.NewOperationContext(http.MethodGet, "/v1/books")
	if err != nil {
		return err
	}
	booksListOp.AddSecurity(securityName)
	booksListOp.SetSummary("Retrieve list of books")
	booksListOp.SetDescription("Retrieve list of books")
	booksListOp.SetID("booksList")
	booksListOp.AddReqStructure(domain.BookListParams{})
	booksListOp.AddRespStructure([]domain.BookResponse{})
	booksListOp.AddRespStructure(failure.Error{}, openapi.WithHTTPStatus(http.StatusBadRequest))
	booksListOp.SetTags("Books")

	err = reflector.AddOperation(booksListOp)
	if err != nil {
		return err
	}

	// BookPutFun
	bookPutOp, err := reflector.NewOperationContext(http.MethodPut, "/v1/books")
	if err != nil {
		return err
	}
	bookPutOp.AddSecurity(securityName)
	bookPutOp.SetSummary("Edit a book")
	bookPutOp.SetDescription("Edit a book")
	bookPutOp.SetID("bookEdit")
	bookPutOp.AddReqStructure(domain.EditBookBody{})
	bookPutOp.SetTags("Books")

	bookPutOp.AddRespStructure(nil, openapi.WithHTTPStatus(http.StatusNoContent))
	bookPutOp.AddRespStructure(failure.Error{}, openapi.WithHTTPStatus(http.StatusBadRequest))
	bookPutOp.AddRespStructure(failure.Error{}, openapi.WithHTTPStatus(http.StatusInternalServerError))

	err = reflector.AddOperation(bookPutOp)
	if err != nil {
		return err
	}

	// BookSavePostFun
	bookSaveOp, err := reflector.NewOperationContext(http.MethodPost, "/v1/books")
	if err != nil {
		return err
	}
	bookSaveOp.AddSecurity(securityName)
	bookSaveOp.SetSummary("Save an already existing book")
	bookSaveOp.SetDescription("Save an already existing book")
	bookSaveOp.SetID("bookSave")
	bookSaveOp.AddReqStructure(domain.SaveBookBody{})
	bookSaveOp.AddRespStructure(domain.BookResponse{})
	bookSaveOp.AddRespStructure(failure.Error{}, openapi.WithHTTPStatus(http.StatusBadRequest))
	bookSaveOp.AddRespStructure(failure.Error{}, openapi.WithHTTPStatus(http.StatusInternalServerError))
	bookSaveOp.SetTags("Books")

	err = reflector.AddOperation(bookSaveOp)
	if err != nil {
		return err
	}

	// BookDeleteFun
	bookDeleteOp, err := reflector.NewOperationContext(http.MethodDelete, "/v1/books/{bookUUID}")
	if err != nil {
		return err
	}
	bookDeleteOp.AddSecurity(securityName)
	bookDeleteOp.SetSummary("Delete a book")
	bookDeleteOp.SetDescription("Delete a book")
	bookDeleteOp.SetID("bookDelete")
	// bookDeleteOp.AddReqStructure(domain.BookGetPath{})
	bookDeleteOp.AddRespStructure(nil, openapi.WithHTTPStatus(http.StatusNoContent))
	bookDeleteOp.AddRespStructure(failure.Error{}, openapi.WithHTTPStatus(http.StatusBadRequest))
	bookDeleteOp.AddRespStructure(failure.Error{}, openapi.WithHTTPStatus(http.StatusInternalServerError))
	bookDeleteOp.SetTags("Books")

	err = reflector.AddOperation(bookDeleteOp)
	if err != nil {
		return err
	}
	return nil
}
