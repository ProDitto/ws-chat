/**
 * @typedef {object} User
 * @property {number} id
 * @property {string} username
 * @property {string} display_name
 * @property {string|null} profile_image_url
 * @property {string} created_at
 * @property {string} updated_at
 * @property {string|null} last_seen_at
 */

/**
 * @typedef {object} Message
 * @property {number} id
 * @property {number} conversation_id
 * @property {number} sender_id
 * @property {'text' | 'media'} type
 * @property {string} content
 * @property {string} created_at
 * @property {User} [sender]
 */

