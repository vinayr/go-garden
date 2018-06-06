TARGET	= garden
LDFLAGS	= -ldflags='-s -w'

build:
		@go build $(LDFLAGS) -o $(TARGET)

clean:
		@rm -f $(TARGET)

run: build
		@godotenv -f .env.dev ./$(TARGET)
