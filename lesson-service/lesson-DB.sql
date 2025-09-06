CREATE TABLE units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    section_count INT,
    order_index INT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    type TEXT,
    lesson_count INT,
    order_index INT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    section_id UUID NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    title TEXT,
    description TEXT,
    exercise_count INT,
    order_index INT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    subtype TEXT,
    instruction TEXT,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    current_unit UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE
    current_section UUID NOT NULL REFERENCES sections(id) ON DELETE CASCADE
    current_lesson UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE
);

-- Lessons
INSERT INTO lessons (unit_id, section_id, title, description, exercise_count, order_index)
VALUES
((SELECT id FROM units WHERE title = 'Greetings'), (SELECT id FROM sections WHERE type = 'core' AND unit_id = (SELECT id FROM units WHERE title = 'Greetings')), 
 'Lesson 1: Basic Greetings', 'Learn basic greetings, enough to get a conversation going', 3, 1),
((SELECT id FROM units WHERE title = 'Greetings'), (SELECT id FROM sections WHERE type = 'core' AND unit_id = (SELECT id FROM units WHERE title = 'Greetings')), 
 'Lesson 2: How are you?', 'Learn simple conversational phrases, that build on the basic greetings to create natural connections', 3, 2),
((SELECT id FROM units WHERE title = 'Food'), (SELECT id FROM sections WHERE type = 'core' AND unit_id = (SELECT id FROM units WHERE title = 'Food')), 
 'Lesson 1: Fruits', 'Learn fruit vocabulary', 3, 1),
((SELECT id FROM units WHERE title = 'Food'), (SELECT id FROM sections WHERE type = 'core' AND unit_id = (SELECT id FROM units WHERE title = 'Food')), 
 'Lesson 2: Meals', 'Learn meal-related vocabulary', 3, 2);

-- Exercises
INSERT INTO exercises (lesson_id, type, subtype, instruction, data)
VALUES
-- Exercises for Lesson 1: Basic Greetings
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'translation', 'block_build', 'Build the Amharic translation from the blocks.',
  '{"prompt_text": "Hello, how are you?", "prompt_audio_url": "hello_how_are_you.mp3", "blocks": ["እንዴት", "ነህ?", "ሰላም"], "correct_answer": "ሰላም እንዴት ነህ?"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"prompt_text": "Thank you", "prompt_audio_url": "thank_you.mp3", "correct_answer": "አመሰግናለሁ"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"prompt_text": "Good morning", "prompt_audio_url": "good_morning.mp3", "correct_answer": "እንደምን አደርክ"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'complete_sentence', 'partial_free_text', 'Complete the sentence.',
  '{"reference_text": "I am fine.", "display_text": "እኔ ____ ነኝ።", "correct_answer": "ደህና"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Basic Greetings'), 'speaking', NULL, 'Speak this sentence aloud.',
  '{"prompt_text": "Goodbye", "prompt_audio_url": "goodbye.mp3", "correct_answer": "ደህና ሁን"}'),

-- Exercises for Lesson 2: How are you?
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'translation', 'block_build', 'Build the Amharic translation from the blocks.',
  '{"prompt_text": "I am fine, thank you.", "prompt_audio_url": "i_am_fine_thank_you.mp3", "blocks": ["ነኝ", "አመሰግናለሁ", "ደህና", "፣"], "correct_answer": "ደህና ነኝ፣ አመሰግናለሁ"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"prompt_text": "What is your name?", "prompt_audio_url": "what_is_your_name.mp3", "correct_answer": "ስምህ ማን ነው?"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'fill_in_blank', NULL, 'Choose the correct word for the blank.',
  '{"display_text": "ስሜ ____ ነው።", "options": [{"id": 0, "text": "ማርያም"}, {"id": 1, "text": "ዮሐንስ"}, {"id": 2, "text": "ሰላም"}], "correct_option_id": 0}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'complete_sentence', 'partial_free_text', 'Complete the sentence.',
  '{"reference_text": "My name is John.", "display_text": "ስሜ ____ ነው።", "correct_answer": "ዮሐንስ"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: How are you?'), 'speaking', NULL, 'Speak this sentence aloud.',
  '{"prompt_text": "How are you?", "prompt_audio_url": "how_are_you.mp3", "correct_answer": "እንዴት ነሽ?"}'),

-- Exercises for Lesson 1: Fruits
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"prompt_text": "Apple", "prompt_audio_url": "apple.mp3", "correct_answer": "ፖም"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'fill_in_blank', NULL, 'Choose the correct word for the blank.',
  '{"display_text": "እኔ ____ እበላለሁ።", "options": [{"id": 0, "text": "ፖም"}, {"id": 1, "text": "ውሻ"}, {"id": 2, "text": "መኪና"}], "correct_option_id": 0}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'translation', 'block_build', 'Build the Amharic translation from the blocks.',
  '{"prompt_text": "I eat a banana.", "prompt_audio_url": "i_eat_a_banana.mp3", "blocks": ["ሙዝ", "እበላለሁ", "እኔ"], "correct_answer": "እኔ ሙዝ እበላለሁ።"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"prompt_text": "Orange", "prompt_audio_url": "orange.mp3", "correct_answer": "ብርቱካን"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 1: Fruits'), 'fill_in_blank', NULL, 'Choose the correct word for the blank.',
  '{"display_text": "ይህ ____ ነው።", "options": [{"id": 0, "text": "ሙዝ"}, {"id": 1, "text": "ፖም"}, {"id": 2, "text": "ብርቱካን"}], "correct_option_id": 0}'),

-- Exercises for Lesson 2: Meals
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"prompt_text": "Breakfast", "prompt_audio_url": "breakfast.mp3", "correct_answer": "ቁርስ"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'translation', 'block_build', 'Build the Amharic translation from the blocks.',
  '{"prompt_text": "I drink water.", "prompt_audio_url": "i_drink_water.mp3", "blocks": ["እጠጣለሁ", "ውሃ", "እኔ"], "correct_answer": "እኔ ውሃ እጠጣለሁ።"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'fill_in_blank', NULL, 'Choose the correct word for the blank.',
  '{"display_text": "ምግቡ ____ ነው።", "options": [{"id": 0, "text": "ጣፋጭ"}, {"id": 1, "text": "ቀዝቃዛ"}, {"id": 2, "text": "ትኩስ"}], "correct_option_id": 0}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'translation', 'free_text', 'Type the Amharic translation.',
  '{"prompt_text": "Lunch", "prompt_audio_url": "lunch.mp3", "correct_answer": "ምሳ"}'),
((SELECT id FROM lessons WHERE title = 'Lesson 2: Meals'), 'speaking', NULL, 'Speak this sentence aloud.',
  '{"prompt_text": "Dinner", "prompt_audio_url": "dinner.mp3", "correct_answer": "እራት"}');



INSERT INTO user_progress (user_id, current_unit, current_section, current_lesson)
VALUES
('5c8530b5-5056-466c-b82f-82874cf226c4', 
'a5fd6809-9e7e-4f4b-8565-f6c1e7b08b8a', 
'cccf142a-e056-4757-a437-4b93e845231b',
'31eca335-5566-4d4d-bf88-c908cef3e3c7'
)