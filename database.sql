CREATE DATABASE IF NOT EXISTS test_db
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE test_db;

SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS restaurant_staff;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS queues;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS `tables`;
DROP TABLE IF EXISTS restaurant_images;
DROP TABLE IF EXISTS restaurants;
DROP TABLE IF EXISTS users;
SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE USERS (
  user_id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  phone VARCHAR(20),
  password VARCHAR(255) NOT NULL,
  role ENUM('customer','owner','staff') DEFAULT 'customer',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ROLES (
  role_id INT AUTO_INCREMENT PRIMARY KEY,
  role_name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS RESTAURANTS (
  restaurant_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  owner_id INT NOT NULL,
  name VARCHAR(150) NOT NULL,
  description TEXT,
  location VARCHAR(255),
  phone VARCHAR(20),
  open_time TIME,
  close_time TIME,
  status ENUM('active', 'inactive', 'suspended') NOT NULL DEFAULT 'active',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (owner_id) REFERENCES USERS(user_id)
    ON DELETE CASCADE
);

CREATE TABLE RESTAURANT_STAFF (
  staff_id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  restaurant_id INT NOT NULL,
  role_id INT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  UNIQUE (user_id, restaurant_id),

  FOREIGN KEY (user_id) REFERENCES USERS(user_id) ON DELETE CASCADE,
  FOREIGN KEY (restaurant_id) REFERENCES RESTAURANTS(restaurant_id) ON DELETE CASCADE,
  FOREIGN KEY (role_id) REFERENCES ROLES(role_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS RESTAURANT_IMAGES (
  image_id INT AUTO_INCREMENT PRIMARY KEY,
  restaurant_id INT NOT NULL,
  image_url TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (restaurant_id) REFERENCES RESTAURANTS(restaurant_id)
    ON DELETE CASCADE
);

CREATE TABLE `TABLES` (
  table_id INT AUTO_INCREMENT PRIMARY KEY,
  restaurant_id INT NOT NULL,
  table_number INT NOT NULL,
  capacity INT NOT NULL,
  status ENUM('available','occupied','reserved') DEFAULT 'available',
  UNIQUE (restaurant_id, table_number),
  FOREIGN KEY (restaurant_id) REFERENCES RESTAURANTS(restaurant_id) ON DELETE CASCADE
);


CREATE TABLE BOOKINGS (
  booking_id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  restaurant_id INT NOT NULL,
  table_id INT,

  booking_date DATE NOT NULL,
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,

  number_of_people INT NOT NULL,
  total_price DECIMAL(10,2) DEFAULT 0,

  status ENUM('pending', 'confirmed', 'completed', 'cancelled') DEFAULT 'pending',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (table_id, booking_date, start_time),
  FOREIGN KEY (user_id) REFERENCES USERS(user_id) ON DELETE CASCADE,
  FOREIGN KEY (restaurant_id) REFERENCES RESTAURANTS(restaurant_id) ON DELETE CASCADE,
  FOREIGN KEY (table_id) REFERENCES TABLES(table_id) ON DELETE SET NULL
);


CREATE TABLE QUEUES (
  queue_id INT AUTO_INCREMENT PRIMARY KEY,
  restaurant_id INT NOT NULL,
  user_id INT NOT NULL,
  queue_number INT NOT NULL,
  number_of_people INT NOT NULL,
  status ENUM('waiting', 'calling', 'completed', 'cancelled') DEFAULT 'waiting',
  UNIQUE (restaurant_id, queue_number),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (restaurant_id) REFERENCES RESTAURANTS(restaurant_id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES USERS(user_id) ON DELETE CASCADE
);


CREATE TABLE REVIEWS (
 review_id INT AUTO_INCREMENT PRIMARY KEY,
 user_id INT NOT NULL,
 restaurant_id INT NOT NULL,
 rating TINYINT CHECK (rating BETWEEN 1 AND 5),
 comment TEXT,
 created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
 FOREIGN KEY (user_id) REFERENCES USERS(user_id) ON DELETE CASCADE,
 FOREIGN KEY (restaurant_id) REFERENCES RESTAURANTS(restaurant_id) ON DELETE CASCADE
);

CREATE TABLE NOTIFICATIONS (
  notification_id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  message TEXT NOT NULL,
  type ENUM('booking','queue') NOT NULL,
  status ENUM('unread','read') DEFAULT 'unread',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES USERS(user_id)
    ON DELETE CASCADE
);

CREATE INDEX idx_restaurant_name ON RESTAURANTS(name);
CREATE INDEX idx_restaurant_location ON RESTAURANTS(location);

-- เพิ่มลูกค้าและเจ้าของร้าน
INSERT INTO USERS (name, email, phone, password) VALUES 
('คุณลูกค้า ใจดี', 'jane@email.com', '0812345678', '123456'),
('คุณลูกค้า ขาประจำ', 'boy@email.com', '0811445978', '123456'),
('คุณลูกค้า สมหญิง', 'somgirl@email.com', '0682325171', '123456'),
('เจ้าของร้าน สมชาย', 'owner1@email.com', '0916345609', '123456'),
('เจ้าของร้าน สาขาหนึ่ง', 'number_one@email.com', '0814325678', '123456'),
('เจ้าของร้าน สมปอง', 'pong@email.com', '0815432678', '123456');

INSERT INTO ROLES (role_name) VALUES
('owner'),
('manager'),
('staff');

-- เพิ่มร้าน
INSERT INTO RESTAURANTS (owner_id, name, location, phone, open_time, close_time) VALUES 
(4, 'ร้านข้าวแกงรังสิต', 'ตึก 15 ม.รังสิต', '0811111111', '08:00:00', '16:00:00'),
(5, 'Sushi House', 'ฟิวเจอร์พาร์ครังสิต', '0822222222', '10:00:00', '22:00:00'),
(6, 'Cafe Chill Bangkhen', 'บางเขน กรุงเทพฯ', '0833333333', '09:00:00', '18:00:00');

INSERT INTO RESTAURANT_STAFF (user_id, restaurant_id, role_id) VALUES
(4, 1, 1), -- owner ร้าน 1
(5, 2, 1), -- owner ร้าน 2
(6, 3, 1); -- owner ร้าน 3

INSERT INTO TABLES (restaurant_id, table_number, capacity, status) VALUES 
(1, 1, 2, 'available'),
(1, 2, 4, 'occupied');

-- เพิ่มการจองและคิว
INSERT INTO BOOKINGS (user_id, restaurant_id, table_id, booking_date, start_time, end_time, number_of_people, total_price, status)
VALUES (1, 1, 1, '2026-05-10', '12:00:00', '14:00:00', 2, 500, 'confirmed');

INSERT INTO QUEUES (restaurant_id, user_id, queue_number, number_of_people, status) VALUES 
(1, 2, 1, 3, 'waiting');

INSERT INTO REVIEWS (user_id, restaurant_id,rating, comment) VALUES 
(1, 1, 5, 'อร่อยมาก!'),
(2, 1, 4, 'บรรยากาศดี'),
(3, 1, 2, 'อาหารเสิร์ฟช้า');

INSERT INTO NOTIFICATIONS (user_id, message, type) VALUES 
(1, 'Your booking is confirmed', 'booking'),
(2, 'Your booking is deny', 'booking'),
(3, 'Your booking is confirmed', 'booking');


-- TEST QUERY
SELECT * FROM USERS;
SELECT * FROM RESTAURANTS;
SELECT * FROM BOOKINGS; 
SELECT * FROM QUEUES;

-- OPTIONAL FEATURES
SELECT * FROM REVIEWS;
SELECT * FROM NOTIFICATIONS;