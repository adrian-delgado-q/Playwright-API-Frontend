import { test, expect } from '@playwright/test';

const API_BASE_URL = 'http://localhost:8080/api/v1';

test.describe('Books API E2E Tests @api @e2e', () => {
  test('should get health status @smoke', async ({ request }) => {
    const response = await request.get('http://localhost:8080/health');
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.status).toBe('ok');
  });

  test('should get all books @smoke @books', async ({ request }) => {
    const response = await request.get(`${API_BASE_URL}/books`);
    expect(response.ok()).toBeTruthy();

    const books = await response.json();
    expect(Array.isArray(books)).toBeTruthy();
    expect(books.length).toBeGreaterThan(0);

    // Verify book structure
    const firstBook = books[0];
    expect(firstBook).toHaveProperty('id');
    expect(firstBook).toHaveProperty('title');
    expect(firstBook).toHaveProperty('author');
    expect(firstBook).toHaveProperty('isbn');
    expect(firstBook).toHaveProperty('year');
  });

  test('should get a specific book by ID @books', async ({ request }) => {
    // First get all books to find a valid ID
    const booksResponse = await request.get(`${API_BASE_URL}/books`);
    expect(booksResponse.ok()).toBeTruthy();

    const books = await booksResponse.json();
    expect(books.length).toBeGreaterThan(0);

    const firstBook = books[0];

    // Now get that specific book
    const response = await request.get(`${API_BASE_URL}/books/${firstBook.id}`);
    expect(response.ok()).toBeTruthy();

    const book = await response.json();
    expect(book.id).toBe(firstBook.id);
    expect(book.title).toBe(firstBook.title);
    expect(book.author).toBe(firstBook.author);
    expect(book.isbn).toBe(firstBook.isbn);
  });

  test('should return 404 for non-existent book @books', async ({
    request,
  }) => {
    const response = await request.get(`${API_BASE_URL}/books/999`);
    expect(response.status()).toBe(404);
  });

  test('should create a new book @books @crud', async ({ request }) => {
    const uniqueISBN = `978${Date.now().toString().slice(-9)}`;
    const newBook = {
      title: 'Test Driven Development',
      author: 'Kent Beck',
      isbn: uniqueISBN,
      year: 2002,
    };

    const response = await request.post(`${API_BASE_URL}/books`, {
      data: newBook,
    });

    expect(response.status()).toBe(201);

    const createdBook = await response.json();
    expect(createdBook.id).toBeDefined();
    expect(createdBook.title).toBe(newBook.title);
    expect(createdBook.author).toBe(newBook.author);
    expect(createdBook.isbn).toBe(newBook.isbn);
    expect(createdBook.year).toBe(newBook.year);

    // Verify book was actually created by fetching it
    const getResponse = await request.get(
      `${API_BASE_URL}/books/${createdBook.id}`
    );
    expect(getResponse.ok()).toBeTruthy();

    const fetchedBook = await getResponse.json();
    expect(fetchedBook.title).toBe(newBook.title);
  });

  test('should validate required fields when creating book @books @validation', async ({
    request,
  }) => {
    const invalidBook = {
      title: 'Incomplete Book',
      // Missing author and isbn
    };

    const response = await request.post(`${API_BASE_URL}/books`, {
      data: invalidBook,
    });

    expect(response.status()).toBe(400);
  });

  test('should update an existing book @books @crud', async ({ request }) => {
    // First create a book
    const uniqueISBN = `978${Date.now().toString().slice(-9)}`;
    const newBook = {
      title: 'Original Title',
      author: 'Original Author',
      isbn: uniqueISBN,
      year: 2020,
    };

    const createResponse = await request.post(`${API_BASE_URL}/books`, {
      data: newBook,
    });
    const createdBook = await createResponse.json();

    // Update the book
    const updateData = {
      title: 'Updated Title',
      year: 2021,
    };

    const updateResponse = await request.put(
      `${API_BASE_URL}/books/${createdBook.id}`,
      {
        data: updateData,
      }
    );

    expect(updateResponse.ok()).toBeTruthy();

    const updatedBook = await updateResponse.json();
    expect(updatedBook.title).toBe('Updated Title');
    expect(updatedBook.author).toBe('Original Author'); // Should remain unchanged
    expect(updatedBook.year).toBe(2021);
  });

  test('should delete a book @books @crud', async ({ request }) => {
    // First create a book
    const uniqueISBN = `978${Date.now().toString().slice(-9)}`;
    const newBook = {
      title: 'Book to Delete',
      author: 'Delete Author',
      isbn: uniqueISBN,
      year: 2023,
    };

    const createResponse = await request.post(`${API_BASE_URL}/books`, {
      data: newBook,
    });
    const createdBook = await createResponse.json();

    // Delete the book
    const deleteResponse = await request.delete(
      `${API_BASE_URL}/books/${createdBook.id}`
    );
    expect(deleteResponse.status()).toBe(204);

    // Verify book is deleted
    const getResponse = await request.get(
      `${API_BASE_URL}/books/${createdBook.id}`
    );
    expect(getResponse.status()).toBe(404);
  });

  test('should handle concurrent requests @books @performance', async ({
    request,
  }) => {
    // Test concurrent GET requests
    const requests = Array(5)
      .fill(null)
      .map(() => request.get(`${API_BASE_URL}/books`));

    const responses = await Promise.all(requests);

    responses.forEach(response => {
      expect(response.ok()).toBeTruthy();
    });
  });

  test('should return proper CORS headers @books @cors', async ({
    request,
  }) => {
    const response = await request.get(`${API_BASE_URL}/books`);
    expect(response.ok()).toBeTruthy();

    const headers = response.headers();
    expect(headers['access-control-allow-origin']).toBe('*');
  });

  test('should handle OPTIONS request for CORS @books @cors', async ({
    request,
  }) => {
    const response = await request.fetch(`${API_BASE_URL}/books`, {
      method: 'OPTIONS',
    });

    expect(response.status()).toBe(200);

    const headers = response.headers();
    expect(headers['access-control-allow-origin']).toBe('*');
    expect(headers['access-control-allow-methods']).toContain('GET');
    expect(headers['access-control-allow-methods']).toContain('POST');
    expect(headers['access-control-allow-methods']).toContain('PUT');
    expect(headers['access-control-allow-methods']).toContain('DELETE');
  });

  test('should maintain data consistency across operations @books @integration', async ({
    request,
  }) => {
    // Create a unique book to avoid conflicts with other tests
    const uniqueISBN = `978${Date.now().toString().slice(-9)}`;
    const newBook = {
      title: 'Consistency Test Book',
      author: 'Test Author',
      isbn: uniqueISBN,
      year: 2024,
    };

    // Create a new book
    const createResponse = await request.post(`${API_BASE_URL}/books`, {
      data: newBook,
    });
    const createdBook = await createResponse.json();
    expect(createResponse.status()).toBe(201);

    // Verify the book exists
    const getResponse = await request.get(
      `${API_BASE_URL}/books/${createdBook.id}`
    );
    expect(getResponse.ok()).toBeTruthy();
    const fetchedBook = await getResponse.json();
    expect(fetchedBook.title).toBe(newBook.title);

    // Update the book
    const updateResponse = await request.put(
      `${API_BASE_URL}/books/${createdBook.id}`,
      {
        data: { title: 'Updated Consistency Test Book' },
      }
    );
    expect(updateResponse.ok()).toBeTruthy();

    // Verify the update
    const updatedGetResponse = await request.get(
      `${API_BASE_URL}/books/${createdBook.id}`
    );
    const updatedBook = await updatedGetResponse.json();
    expect(updatedBook.title).toBe('Updated Consistency Test Book');

    // Delete the book
    const deleteResponse = await request.delete(
      `${API_BASE_URL}/books/${createdBook.id}`
    );
    expect(deleteResponse.status()).toBe(204);

    // Verify book is deleted
    const deletedGetResponse = await request.get(
      `${API_BASE_URL}/books/${createdBook.id}`
    );
    expect(deletedGetResponse.status()).toBe(404);
  });
});
