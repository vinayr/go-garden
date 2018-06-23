TARGET	= garden
LDFLAGS	= -ldflags='-s -w'

build:
		@go build $(LDFLAGS) -o $(TARGET)

clean:
		@rm -f $(TARGET)

run: build
		./$(TARGET)
