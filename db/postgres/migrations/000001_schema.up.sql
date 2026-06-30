CREATE TABLE Documents (
    Id BIGSERIAL PRIMARY KEY,
    Guid VARCHAR(40) NOT NULL,
    
    Title VARCHAR(50) NOT NULL,
    Description TEXT NOT NULL,
    Content TEXT NOT NULL,

    Created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Authors (
    Id BIGSERIAL PRIMARY KEY,
    Name TEXT UNIQUE NOT NULL
);

CREATE TABLE Tags (
    Id BIGSERIAL PRIMARY KEY,
    Name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE Links (
    Id BIGSERIAL PRIMARY KEY,
    Link TEXT UNIQUE NOT NULL
);

CREATE TABLE Document_Authors (
    DocumentId BIGINT NOT NULL, 
    AuthorId BIGINT NOT NULL,

    CONSTRAINT pk_DocAuthor PRIMARY KEY (DocumentId, AuthorId),
    CONSTRAINT fk_Document FOREIGN KEY (DocumentId) REFERENCES Documents(Id) ON DELETE CASCADE,
    CONSTRAINT fk_Author FOREIGN KEY (AuthorId) REFERENCES Authors(Id) ON DELETE CASCADE
);

CREATE TABLE Document_Tags (
    DocumentId BIGINT NOT NULL, 
    TagId BIGINT NOT NULL,

    CONSTRAINT pk_DocTag PRIMARY KEY (DocumentId, TagId),
    CONSTRAINT fk_Document FOREIGN KEY (DocumentId) REFERENCES Documents(Id) ON DELETE CASCADE,
    CONSTRAINT fk_Tag FOREIGN KEY (TagId) REFERENCES Tags(Id) ON DELETE CASCADE
);

CREATE TABLE Document_Links (
    DocumentId BIGINT NOT NULL, 
    LinkId BIGINT NOT NULL,

    CONSTRAINT pk_DocLink PRIMARY KEY (DocumentId, LinkId),
    CONSTRAINT fk_Document FOREIGN KEY (DocumentId) REFERENCES Documents(Id) ON DELETE CASCADE,
    CONSTRAINT fk_Link FOREIGN KEY (LinkId) REFERENCES Links(Id) ON DELETE CASCADE
);

CREATE INDEX idx_guid ON Documents(Guid);
CREATE INDEX idx_author ON Authors(Name);
CREATE INDEX idx_tag ON Tags(Name);
CREATE INDEX idx_link ON Links(Link);

CREATE INDEX idx_doc_authors_author_id ON Document_Authors(AuthorId);
CREATE INDEX idx_doc_tags_tag_id ON Document_Tags(TagId);