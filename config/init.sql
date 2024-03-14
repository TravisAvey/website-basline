create table posts (
  id serial not null unique,
  title varchar(64),
  dateposted date,
  dateupdated date,
  imageurl varchar(128),
  content text,
  summary text,
  keywords text,
  slug varchar(64) not null unique,
  primary key(id)
);
create table categories (
  id serial not null unique,
  category varchar(32) not null unique,
  primary key (id)
);
-- bridge table between post and category
create table post_categories (
  post_id int not null references posts,
  category_id int not null references categories,
  primary key (post_id, category_id)
);

-- talbe to hold images for a gallery (maybe other parts of the website too?)
create table photos (
  id serial not null unique,
  image_url varchar(255) not null unique,  -- url endpoint where the image is located
  title varchar(50),                       -- optional title for the image, usuful in a gallery when opened full
  summary varchar(255),                    -- optional summary/description for SEO, gallery, and alt text in img tag
  is_gallery boolean,                      -- if the photo should be used in the gallery or for other parts of website
  primary key (id)

);

-- table to hold categories for images
create table photo_categories (
  id serial not null unique,
  category varchar(32) not null unique,
  primary key(id)
);
-- bridge table between photo and photo_cateogry
create table gallery_categories (
  image_id int not null references photos,
  category_id int not null references photo_categories,
  primary key (image_id, category_id)
);

create table messages (
  id serial not null unique,
  type varchar(10) not null,
  header varchar(255) not null,
  message text not null,
  sent date not null,
  read boolean not null,
  primary key (id)
);

-- users: ref user_role
create table users (
  id serial not null unique,
  name varchar(28) not null unique,
  email varchar(50) not null unique,
  password text not null,
  primary key (id)
);

insert into users(name, email, password) values
  ('Jim Halpert','jim@dunder.com','pampampam'),
  ('Dwight Schrute','thebest@dunder.com','bearsbeatsbattlestar'),
  ('Michael Scott','bestboss@dunder.com','thatswhatshesaid');

insert into posts(title, dateposted, dateupdated, imageurl, content, summary, keywords, slug)
values
  ('Hello World','2023-11-23',NULL,E'https://imgur.com/uzdpuEJ.jpg',E'<h3>Hello World</h3>','Hello World','hello world','hello-world'),
  ('The Second Post','2023-12-07',NULL,E'https://imgur.com/uzdpuEJ.jpg',E'<h3>A Post</h3>','This is a post','test','second-post'),
  ('The Post','2023-12-14',NULL,E'https://imgur.com/uzdpuEJ.jpg',E'<h3>The Post</h3>','The third post','test','more-tests'),
  ('Lets Have another Post','2024-01-14',NULL,E'https://imgur.com/uzdpuEJ.jpg',E'<h3>A simple test, hmmm</h3>','May the Forth test - test','test','another-test');

insert into categories(category) values
  ('general'),
  ('off topic'),
  ('news'),
  ('tutorial'),
  ('announcement');

insert into post_categories(post_id, category_id) values
  (1, (select id from categories where id=1) ),
  (1, (select id from categories where id=2) ),
  (2, (select id from categories where id=3) ),
  (3, (select id from categories where id=5) ),
  (4, (select id from categories where id=4) );

insert into photo_categories(category) values
  ('portrait'),
  ('landscape'),
  ('headshot'),
  ('abstract'),
  ('wedding'),
  ('family');

insert into photos(image_url, title, summary, is_gallery)
  values 
  (
    'https://neonspot-images.nyc3.digitaloceanspaces.com/matterhorn.jpg',
    'The Matterhorn',
    'The Matterhorn at sunrise',
    true
  ),
  (
    'https://picsum.photos/1920/1080',
    'Random Image #1',
    'Lorem Picsum Images',
    true
  ),
  (
    'https://picsum.photos/1000/1000',
    'Random Image #2',
    'Lorem Picsum Images',
    true
  ),
  (
    'https://picsum.photos/1080/1920',
    'Random Image #3',
    'Lorem Picsum Images',
    true
  );

insert into gallery_categories(image_id, category_id) values
  (1, (select id from photo_categories where id=2));


insert into messages(type, header, message, sent, read) values
  ('alert', 'New Subscriber', 'Jim has joined! jim@dunder.com', '2020-05-31', FALSE),
  ('alert', 'New Message', 'A test message from the database', '2020-01-18', FALSE),
  ('alert', 'New Subscriber', 'postgresql has joined! email: psql@bash', '2020-08-24', FALSE),
  ('contact', 'Jim Halpert;j.halpert@dunder.com', 'This is Jim Halpert.', '2020-03-18', FALSE),
  ('contact', 'Dwight Schrute;d.schrute@dudner.com', 'Bears eat beats', '2020-05-14', FALSE);
