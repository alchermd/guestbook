USE guestbook;

CREATE TABLE messages( 
	id INT AUTO_INCREMENT PRIMARY KEY, 
	name VARCHAR(50), 
	message TEXT NOT NULL, 
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
