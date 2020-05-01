/* Migration to add IsAnonymous and GradientColor fields to the participant
 */
ALTER TABLE participants ADD COLUMN is_anonymous boolean NOT NULL DEFAULT True;
ALTER TABLE participants ADD COLUMN gradient_color varchar(36);

